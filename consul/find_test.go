package consul

import "testing"

func TestFindService(t *testing.T) {
	RegisterWithoutTLS(8080, "192.168.1.225:8500")
	FindService()
	for k, v := range Services.Services {
		t.Logf("%v: %v", k, v)
	}
}
