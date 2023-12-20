package main

import (
	"log"
	"net"
	"time"

	capi "github.com/hashicorp/consul/api"
)

const (
	ttl     = time.Second * 8
	CheckID = "check_health"
)

type Proxy struct {
	cli           *capi.Client
	containerName string
}

// Get a new client
func (P *Proxy) NewProxy(cn string) *Proxy {
	client, err := capi.NewClient(capi.DefaultConfig())
	if err != nil {
		panic(err)
	}
	return &Proxy{
		cli:           client,
		containerName: cn,
	}

}

func (P *Proxy) Start() error {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
		return err
	}
	P.registerService()
	go P.updatehealthcheck()
	P.acceptloop(lis)
	return nil
}

func (P *Proxy) acceptloop(ln net.Listener) {
	for {
		_, err := ln.Accept()
		if err != nil {
			log.Fatal(err)

		}
	}
}
func (P *Proxy) updatehealthcheck() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		if err := P.cli.Agent().UpdateTTL(CheckID, "online", capi.HealthPassing); err != nil {
			log.Fatal(err)
		}

		<-ticker.C
	}

}
func (P *Proxy) registerService() {
	check := &capi.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        CheckID,
	}
	register := &capi.AgentServiceRegistration{
		ID:      "login_service",
		Name:    P.containerName,
		Tags:    []string{"login"},
		Address: "127.0.0.1",
		Port:    3000,
		Check:   check,
	}
	if err := P.cli.Agent().ServiceRegister(register); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var P *Proxy
	p := P.NewProxy("cn")
	if err := p.Start(); err != nil {

		log.Fatal(err)

	}

}
