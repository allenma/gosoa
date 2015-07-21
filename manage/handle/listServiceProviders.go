package handle

import (
	"encoding/json"

	"github.com/go-martini/martini"
	"github.com/allenma/gosoa/registry"
)

func ListServiceProviders(params martini.Params, reg registry.Registry) string {
	serviceName := params["name"]
	providers, err := reg.ListServiceProviders(serviceName)
	if err != nil {
		return err.Error()
	}
	jsonRes, _ := json.Marshal(providers)
	return string(jsonRes)
}
