package consul

import (
	"fmt"
	"testing"
)

func Test_getCurrIP(t *testing.T) {
	fmt.Println(getCurrIP())
}

func TestRegisterWithoutTLS(t *testing.T) {
	var port = 8080
	var remote = "192.168.1.225:8500"
	RegisterWithoutTLS(port, remote)
}
