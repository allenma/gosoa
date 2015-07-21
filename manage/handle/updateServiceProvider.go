package handle

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/allenma/gosoa/registry"
	"io/ioutil"
	"net/http"
)

func UpdateServiceProvider(req *http.Request, params martini.Params, reg registry.Registry) string {
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	providerInfo := registry.ProviderInfo{}
	err := json.Unmarshal(body, &providerInfo)
	if err != nil {
		return err.Error()
	}

	serviceName := params["name"]

	err = reg.UpdateServiceProvider(serviceName, providerInfo)
	if err != nil {
		return err.Error()
	}
	return "{\"ok\"}"
}
