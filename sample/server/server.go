package main

import (
	"fmt"
	"github.com/allenma/gosoa/registry"
	"github.com/allenma/gosoa/sample/server/handle"
	"github.com/allenma/gosoa/sample/tutorial"
	//	"time"
	"github.com/allenma/gosoa/export"
)

func NewCalculatorHandler() tutorial.Calculator {
	return new(handle.CalculatorImpl)
}

func main() {
	handler := NewCalculatorHandler()
	processor := tutorial.NewCalculatorProcessor(handler)
	reg := registry.NewRedisRegistry("localhost:6379")
	//	reg,_ := registry.NewZKRegistry([]string{"192.168.148.128:2181"},10*time.Second)
	exporter := export.NewThriftExporter("localhost:9091", reg)
	err := exporter.Export("calculator", processor)
	if err != nil {
		fmt.Println("error when export service", err.Error())
		return
	}

}
