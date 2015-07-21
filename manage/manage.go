// manage.go
package main

import (
	//	"time"
	"github.com/go-martini/martini"
	"github.com/allenma/gosoa/manage/handle"
	"github.com/allenma/gosoa/registry"
)

func main() {
	reg := registry.NewRedisRegistry("localhost:6379")
	//	reg,_ := registry.NewZKRegistry([]string{"192.168.148.128:2181"},10*time.Second)
	m := martini.Classic()
	m.MapTo(reg, (*registry.Registry)(nil))
	m.Get("/service", handle.ListServices)
	m.Get("/service/:name/providers", handle.ListServiceProviders)
	m.Post("/service/:name/updateprovider", handle.UpdateServiceProvider)
	m.RunOnAddr(":8080")
}
