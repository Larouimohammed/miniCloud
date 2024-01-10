package consul

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
	Cli *capi.Client
	// cancel context.CancelFunc
}

// Get a new client
func NewProxy() *ConsulProxy {
	client, err := capi.NewClient(capi.DefaultNonPooledConfig())
	if err != nil {
		log.Println(err)
	}
	return &ConsulProxy{
		Cli: client,
	}

}

var DefaultConsulProxy = NewProxy()

func (P *ConsulProxy) Start(containerName, containerid, ip string, port int) error {
	err := P.registerService(containerName, containerid, ip, port)
	if err != nil {
		return err
	}
	go P.updatehealthcheck(containerName)
	return nil

}

func (P *ConsulProxy) updatehealthcheck(containerName string) error {

	ticker := time.NewTicker(time.Second * 5)
	for {

		<-ticker.C

		if err := P.Cli.Agent().UpdateTTL(containerName, "pass", capi.HealthPassing); err != nil {
			return err
		}

	}

}
func (P *ConsulProxy) registerService(containerName, containerid, ip string, port int) error {

	check := &capi.AgentServiceCheck{
		Name:                           "consul_check",
		DeregisterCriticalServiceAfter: ttl.String(),
		Shell:                          "/bin/bash",
		TLSSkipVerify:                  true,

		// TTL:               ttl.String(),
		CheckID:                containerName,
		DockerContainerID:      containerid,
		FailuresBeforeCritical: 1,
		SuccessBeforePassing:   1,
		Interval:               "10s",
		Args:                   []string{"sh", "-c"},
		// Timeout:  "15s",
	}

	register := &capi.AgentServiceRegistration{
		ID:   containerName,
		Name: containerName + " consul-proxy",
		Tags: []string{containerid},
		// Address: ip,
		// Port:    port,
		Check: check,
	}

	err := P.Cli.Agent().ServiceRegisterOpts(register, capi.ServiceRegisterOpts{ReplaceExistingChecks: true})

	if err != nil {
		log.Printf("Failed to register service: %s:%v with error : %v", ip, port, err)
		return err
	} else {
		log.Printf("successfully register service: %s:%v", ip, port)
		return nil
	}

}
