package command

import (
	"log"
	"time"

	capi "github.com/hashicorp/consul/api"
)

const (
	ttl = time.Second * 8
	// CheckID = "check_health"
)

type ConsulProxy struct {
	cli *capi.Client
	// cancel context.CancelFunc
}

// Get a new client
func (P *ConsulProxy) NewProxy() *ConsulProxy {
	client, err := capi.NewClient(capi.DefaultNonPooledConfig())
	if err != nil {
		log.Println(err)
	}
	return &ConsulProxy{
		cli: client,
	}

}

func (P *ConsulProxy) Start(containerName, containerid, ip string, port int) error {
	P.registerService(containerName, containerid, ip, port)
	// ctx, cancel := context.WithCancel(context.Background())
	// P.cancel = cancel
	go P.updatehealthcheck(containerName)
	return nil
}

// func (P *ConsulProxy) Close() {
// 	P.cancel()

// }

func (P *ConsulProxy) updatehealthcheck(containerName string) {

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:

			if err := P.cli.Agent().UpdateTTL(containerName, "online", capi.HealthPassing); err != nil {
				log.Fatal(err)
			}
			// case <-ctx.Done():
			// 	P.deregisterservice(containerName)

		}

	}

}
func (P *ConsulProxy) registerService(containerName, containerid, ip string, port int) {
	check := &capi.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        containerName,
		DockerContainerID: containerid,
		// Interval:          "5s",
		Timeout: "10s",
		// HTTP:    fmt.Sprintf("http://%s:%v", ip, port),
	}
	register := &capi.AgentServiceRegistration{
		ID:      "Id" + containerName,
		Name:    containerName,
		Tags:    []string{containerid},
		Address: ip,
		Port:    port,
		Check:   check,
	}
	if err := P.cli.Agent().ServiceRegister(register); err != nil {
		log.Fatal(err)
	}

}
func (P *ConsulProxy) Deregisterservice(containerName string) {

	if err := P.cli.Agent().ServiceDeregister(containerName); err != nil {
		log.Println(err)
	}

}
