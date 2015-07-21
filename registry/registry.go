package registry

type Registry interface {
	Register(serviceName string, providerInfo ProviderInfo) error
	Discover(serviceName string, callback RegChangeCallback) ([]ProviderInfo, error)
	ListServices() ([]string, error)
	ListServiceProviders(serviceName string) ([]ProviderInfo, error)
	UpdateServiceProvider(serviceName string, provider ProviderInfo) error
}

type ProviderInfo struct {
	Addr   string `json:"addr"`   //service provider address
	Status int    `json:"status"` //service provider status, 0 is normal, 1 is stopped
	Weight int    `json:"weight"` //the weight of this service provider default is 100
}
type RegChangeCallback func([]ProviderInfo, error)
