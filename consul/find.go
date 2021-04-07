package consul

import (
	consulApi "github.com/hashicorp/consul/api"
	"log"
)

type TServices struct {
	Services map[string]*consulApi.AgentService
	Redo     bool
}

var Services TServices

func FindService() {
	s, err := client.Agent().Services()
	if err != nil {
		log.Fatalf("get services with consul failed, err: %v", err)
	}

	flag := false
	if len(s) != len(Services.Services) {
		flag = true
	}
	for key, server := range s {
		if Services.Services[key] == nil ||
			Services.Services[key].Address != server.Address ||
			Services.Services[key].Port != server.Port {
			flag = true
		}
	}

	if flag {
		Services.Services = s
		Services.Redo = true
	}
}
