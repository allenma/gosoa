// thrift.go
package client

import (
	"errors"
	"reflect"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/allenma/gosoa/registry"
)

type ThriftClient struct {
	Reg              registry.Registry
	ServiceName      string
	LB               LoadBalance
	ServiceClientFun interface{}
	Retry            int
	ProtocolFactory  thrift.TProtocolFactory
	pools            map[string]*ServicePool // transport pools
	providersCache   []registry.ProviderInfo // cache providers in local memory
}

func NewThriftClient(reg registry.Registry, serviceName string, fun interface{}) *ThriftClient {
	tc := &ThriftClient{
		Reg:              reg,
		ServiceName:      serviceName,
		ServiceClientFun: fun,
		LB:               RandomLB(),
		ProtocolFactory:  thrift.NewTBinaryProtocolFactoryDefault(),
		Retry:            0,
		pools:            make(map[string]*ServicePool),
	}
	return tc
}

// return the service implemention and the service provider's address
func (t *ThriftClient) getService() (service interface{}, addr string, err error) {
	if t.providersCache == nil {
		t.providersCache, err = t.Reg.Discover(t.ServiceName, func(providers []registry.ProviderInfo, err error) {
			if err != nil {
				t.providersCache = providers
			}
		})
	}
	if err != nil {
		return
	}
	provider, err := t.LB.Select(t.providersCache)
	if err != nil {
		return
	}
	addr = provider.Addr
	service, err = t.getServiceFromPool(addr)
	return
}

func (t *ThriftClient) Execute(methodName string, args ...interface{}) (res interface{}, err error) {
	for retry := 0; retry <= t.Retry; retry++ {
		res, err = t.doExecute(methodName, args...)
		if err == nil {
			return res, err
		}
	}
	return
}

func (t *ThriftClient) doExecute(methodName string, args ...interface{}) (res interface{}, err error) {
	service, addr, err := t.getService()
	defer t.returnServiceToPool(service, addr)
	if err != nil {
		return
	}

	if err != nil {
		return
	}
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	results := reflect.ValueOf(service).MethodByName(methodName).Call(inputs)

	if len(results) == 0 {
		return nil, nil
	} else if len(results) == 1 {
		if results[0].IsNil() {
			return nil, nil
		}
		return nil, results[0].Interface().(error)
	} else {
		err = findErrorInValues(results)
		return results[0].Interface(), err
	}
}

func findErrorInValues(vals []reflect.Value) error {
	for _, val := range vals {
		if _, ok := val.Interface().(error); ok {
			return val.Interface().(error)
		}
	}
	return nil
}

func (t *ThriftClient) getServiceFromPool(addr string) (service interface{}, err error) {
	servicePool, ok := t.pools[addr]
	if !ok {
		servicePool = newServicePool(addr, t.ServiceClientFun, t.ProtocolFactory, 10, time.Duration(5)*time.Second)
		t.pools[addr] = servicePool
	}
	service, err = servicePool.borrowService()
	return
}

func (t *ThriftClient) returnServiceToPool(service interface{}, addr string) (err error) {
	servicePool, ok := t.pools[addr]
	if ok {
		return servicePool.returnService(service)
	}
	return
}

type ServicePool struct {
	pool             chan interface{}
	serviceCreateFun interface{}
	protocolFactory  thrift.TProtocolFactory
	size             int
	addr             string
	timeout          time.Duration
}

func newServicePool(addr string, serviceCreateFun interface{}, protocolFactory thrift.TProtocolFactory, poolSize int, timeout time.Duration) (pool *ServicePool) {
	pool = new(ServicePool)
	pool.addr = addr
	pool.size = poolSize
	pool.timeout = timeout
	pool.serviceCreateFun = serviceCreateFun
	pool.protocolFactory = protocolFactory

	pool.pool = make(chan interface{}, poolSize)
	for i := 0; i < pool.size; i++ {
		pool.pool <- nil
	}
	return
}

func (s *ServicePool) borrowService() (service interface{}, err error) {
	if s.pool == nil {
		return nil, errors.New("transport pool is not initialized")
	}

	select {
	case service = <-s.pool:
		if service == nil {
			service, err = s.createService(s.addr)
		}
	case <-time.After(s.timeout):
		err = errors.New("timeout when get transport from pool")
	}
	return
}

func (s *ServicePool) createService(addr string) (service interface{}, err error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

	socket, err := thrift.NewTSocket(addr)
	transport := transportFactory.GetTransport(socket)
	err = transport.Open()
	if err != nil {
		return nil, err
	} else {
		rfun := reflect.ValueOf(s.serviceCreateFun)
		result := rfun.Call([]reflect.Value{reflect.ValueOf(transport), reflect.ValueOf(s.protocolFactory)})
		if len(result) != 1 {
			return nil, errors.New("func must return 1 result")
		}
		return result[0].Interface(), err
	}

}

func (s *ServicePool) returnService(service interface{}) (err error) {
	if s.pool == nil {
		return errors.New("service pool is nil")
	}
	s.pool <- service
	return
}
