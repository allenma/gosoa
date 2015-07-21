// loadbalance.go
package client

import (
	"errors"
	"math/rand"
	"time"

	"github.com/allenma/gosoa/registry"
)

type LoadBalance interface {
	Select(infos []registry.ProviderInfo) (registry.ProviderInfo, error)
}

type RandomLoadBalance struct {
	rnd *rand.Rand
}

func RandomLB() *RandomLoadBalance {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &RandomLoadBalance{r}
}

func (r *RandomLoadBalance) Select(providers []registry.ProviderInfo) (provider registry.ProviderInfo, err error) {
	l := len(providers)
	if l <= 0 {
		return registry.ProviderInfo{}, errors.New("empty proivder")
	}
	totalWeight := 0
	availableProviders := make([]registry.ProviderInfo, 0, l)
	for _, provider := range providers {
		if provider.Status == 0 {
			totalWeight += provider.Weight
			availableProviders = append(availableProviders, provider)
		}
	}
	l = len(availableProviders)
	if l <= 0 {
		return registry.ProviderInfo{}, errors.New("no available proivder")
	}
	rndNum := r.rnd.Intn(totalWeight)
	tmpWeight := 0
	i := 0
	for ; i < l; i++ {
		tmpWeight += availableProviders[i].Weight
		if tmpWeight >= rndNum {
			break
		}
	}
	return availableProviders[i], nil
}

type RoundRobinLoadBalance struct {
	availableProviders []registry.ProviderInfo //cached available providers
	maxweight          int            //max weight
	cursor             int            //current selected server
	cw                 int            //current weight
}

func RoundRobinLB() *RoundRobinLoadBalance {
	lb := &RoundRobinLoadBalance{}
	lb.Reset()
	return lb
}
func (r *RoundRobinLoadBalance) Reset() {
	r.availableProviders = nil
	r.cursor = -1
	r.cw = 0
	r.maxweight = -1
}

func (r *RoundRobinLoadBalance) Select(providers []registry.ProviderInfo) (provider registry.ProviderInfo, err error) {
	if r.availableProviders == nil {
		err = r.Initialize(providers)
		if err != nil {
			return
		}
	}
	n := len(r.availableProviders)
	gcd := 1 //greatest common divisor
	for {
		r.cursor = (r.cursor + 1) % n
		if r.cursor == 0 {
			r.cw = r.cw - gcd
			if r.cw <= 0 {
				r.cw = r.maxweight
			}
		}

		if r.availableProviders[r.cursor].Weight >= r.cw {
			return r.availableProviders[r.cursor], nil
		}
	}
	return
}

func (r *RoundRobinLoadBalance) Initialize(providers []registry.ProviderInfo) (err error) {
	l := len(providers)
	if l <= 0 {
		return errors.New("empty addresses")
	}
	r.maxweight = -1
	r.availableProviders = make([]registry.ProviderInfo, 0, l)
	for _, provider := range providers {
		if provider.Status == 0 {
			if provider.Weight > r.maxweight {
				r.maxweight = provider.Weight
			}
			r.availableProviders = append(r.availableProviders, provider)
		}
	}
	l = len(r.availableProviders)
	if l <= 0 {
		return errors.New("no available provider")
	}
	return
}
