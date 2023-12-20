package command

import (
	"log"
	"time"

	capi "github.com/hashicorp/consul/api"
)

const (
	ttl     = time.Second * 8
	// CheckID = "check_health"
)

type ConsulProxy struct {
	cli           *capi.Client

}

// Get a new client
func (P *ConsulProxy) NewProxy() *ConsulProxy {
	client, err := capi.NewClient(capi.DefaultNonPooledConfig())
	if err != nil {
		panic(err)
	}
	return &ConsulProxy{
		cli:           client,
	}

}

func (P *ConsulProxy) Start(containerName string) error {
	P.registerService(containerName)
	go P.updatehealthcheck(containerName)
	return nil
}


func (P *ConsulProxy) updatehealthcheck(containerName string) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		if err := P.cli.Agent().UpdateTTL(containerName, "online", capi.HealthPassing); err != nil {
			log.Fatal(err)
		}

		<-ticker.C
	}

}
func (P *ConsulProxy) registerService(containerName string) {
	check := &capi.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        containerName,
	}
	register := &capi.AgentServiceRegistration{
		ID:      "ID"+containerName,
		Name:    containerName,
		Tags:    []string{"login"},
		Address: "127.0.0.1",
		Port:    3000,
		Check:   check,
	}
	if err := P.cli.Agent().ServiceRegister(register); err != nil {
		log.Fatal(err)
	}
}




