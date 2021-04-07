package exporter

import (
	"github.com/sowhyim/grpc_ui/consul"
	"net/http"
	"os"
	"time"

	"github.com/fullstorydev/grpcui/standalone"
)

type TExporter struct {
	ExportInterval int
	Modify         chan int
	ReFlash        chan bool
	Exit           chan bool
}

var Exporter *TExporter

func init() {
	Exporter = DefaultExporter()
	Exporter.Run()
}

func DefaultExporter() *TExporter {
	return &TExporter{
		ExportInterval: 10,
		Modify:         make(chan int),
		ReFlash:        make(chan bool),
		Exit:           make(chan bool),
	}
}

func (e *TExporter) Run() {
	for {
		select {
		case modify := <-e.Modify:
			e.ExportInterval = modify
		case <-e.ReFlash:
			e.Export()
		case <-e.Exit:
			break
		case <-time.After(time.Second * time.Duration(e.ExportInterval)):
			e.Export()
		}
	}
	os.Exit(0)
}

func (e *TExporter) ModifyInterval(interval int) {
	e.Modify <- interval
}

func (e *TExporter) Export() {
	consul.FindService()
	if consul.Services.Redo {
		var handler []http.Handler
		for i := range consul.Services.Services {
			handler = standalone.Handler()
		}
	}

	http.handl
}
