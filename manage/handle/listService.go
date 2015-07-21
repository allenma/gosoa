package handle

import (
	"encoding/json"
	"github.com/allenma/gosoa/registry"
)

func ListServices(reg registry.Registry) string {
	services, err := reg.ListServices()
	if err != nil {
		return err.Error()
	}
	jsonRes, _ := json.Marshal(services)
	return string(jsonRes)
}
