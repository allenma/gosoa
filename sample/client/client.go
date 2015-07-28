package main

import (
	"fmt"
	//	"time"
	"github.com/allenma/gosoa/client"
	"github.com/allenma/gosoa/registry"
	"github.com/allenma/gosoa/sample/tutorial"
)

func main() {
	reg := registry.NewRedisRegistry("localhost:6379")
	//	reg,_ := registry.NewZKRegistry([]string{"192.168.148.128:2181"},10*time.Second)
	tclient := client.NewThriftClient(reg, "calculator", tutorial.NewCalculatorClientFactory)
	for i := 0; i < 10; i++ {
		var a int32 = 5
		var b int32 = 6
		fmt.Println("time ", i, ":")
		result, err := tclient.ExecuteWithRetry("Add", 5, a, b)

		if err == nil {
			fmt.Println("  add result=", result.(int32))
		} else {
			fmt.Println("  add error:", err.Error())
		}

		_, err = tclient.ExecuteWithRetry("Ping", 5)
		if err == nil {
			fmt.Println("  ping success")
		} else {
			fmt.Println("  ping error:", err.Error())
		}
	}

}
