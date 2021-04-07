package consul

import (
	"fmt"
	"git.chinalbr.com/Server/monitor/define"
	consulApi "github.com/hashicorp/consul/api"
	"log"
	"net"
)

var client *consulApi.Client

func RegisterWithoutTLS(port int, remote string) {
	currIP, ok := getCurrIP()
	if !ok {
		log.Fatal("get current ip failed")
	}

	config := consulApi.DefaultConfig()
	config.Address = remote
	var err error
	client, err = consulApi.NewClient(config)
	if err != nil {
		log.Fatalf("new consul client failed, err: %v", err)
	}

	registration := new(consulApi.AgentServiceRegistration)
	registration.ID = define.ServerID
	registration.Name = define.ServerName
	registration.Port = port
	registration.Tags = []string{define.ServerTag}
	registration.Address = currIP

	check := new(consulApi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
	check.Timeout = "10s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("client registry failed, err: %v", err)
	}

	log.Println("consul client already start!")
}

func getCurrIP() (string, bool) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return "", false
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), true
					}
				}
			}
		}
	}

	return "", false
}
