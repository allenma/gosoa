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
	serviceCreateFun interface{}
	Retry            int
	ProtocolFactory  thrift.TProtocolFactory
	SocketTimeout    time.Duration
	MaxActive        int
	MaxIdle          int
	IdleTimeout      time.Duration
	pools            map[string]*Pool        // transport pools
	providersCache   []registry.ProviderInfo // cache providers in local memory
}

func NewThriftClient(reg registry.Registry, serviceName string, fun interface{}) *ThriftClient {
	tc := &ThriftClient{
		Reg:              reg,
		ServiceName:      serviceName,
		serviceCreateFun: fun,
		LB:               RandomLB(),
		ProtocolFactory:  thrift.NewTBinaryProtocolFactoryDefault(),
		Retry:            0,
		SocketTimeout:    time.Duration(1000) * time.Millisecond,
		IdleTimeout:      time.Duration(1000) * time.Second,
		MaxActive:        8,
		MaxIdle:          8,
		pools:            make(map[string]*Pool),
	}
	return tc
}

// return the service implemention and the service provider's address
func (t *ThriftClient) getService() (service *ServiceWrapper, addr string, err error) {
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
	return t.ExecuteWithRetry(methodName, t.Retry, args...)
}

func (t *ThriftClient) ExecuteWithRetry(methodName string, retry int, args ...interface{}) (res interface{}, err error) {
	for retryTime := 0; retryTime <= retry; retryTime++ {
		res, err = t.doExecute(methodName, args...)
		if err == nil {
			return res, err
		}
	}
	return
}

func (t *ThriftClient) doExecute(methodName string, args ...interface{}) (res interface{}, err error) {
	service, addr, err := t.getService()
	if err != nil {
		return
	}
	defer t.returnServiceToPool(service, addr)

	if err != nil {
		return
	}
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	results := reflect.ValueOf(service.service).MethodByName(methodName).Call(inputs)

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

func (t *ThriftClient) getServiceFromPool(addr string) (service *ServiceWrapper, err error) {
	servicePool, ok := t.pools[addr]
	if !ok {
		servicePool = t.newServicePool(addr)
		t.pools[addr] = servicePool
	}
	conn, err := servicePool.Borrow()
	service, _ = conn.(*ServiceWrapper)
	return service, err
}

func (t *ThriftClient) returnServiceToPool(service *ServiceWrapper, addr string) (err error) {
	servicePool, ok := t.pools[addr]
	if ok {
		return servicePool.Return(service, false)
	}
	return
}

func (t *ThriftClient) newServicePool(addr string) (pool *Pool) {
	pool = &Pool{
		Dial: func() (Conn, error) {
			conn, err := t.createService(addr)
			return conn, err
		},
		Wait:        true,
		MaxActive:   t.MaxActive,
		MaxIdle:     t.MaxIdle,
		IdleTimeout: t.IdleTimeout,
	}
	return pool
}

type ServiceWrapper struct {
	service   interface{}
	transport thrift.TTransport
	err       error
}

func (s *ServiceWrapper) Close() error {
	if s.transport != nil {
		return s.transport.Close()
	}
	return nil
}

func (s *ServiceWrapper) GetErr() error {
	return s.err
}

func (t *ThriftClient) createService(addr string) (service *ServiceWrapper, err error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

	socket, err := thrift.NewTSocket(addr)
	transport := transportFactory.GetTransport(socket)
	err = transport.Open()
	if err != nil {
		return nil, err
	} else {
		rfun := reflect.ValueOf(t.serviceCreateFun)
		result := rfun.Call([]reflect.Value{reflect.ValueOf(transport), reflect.ValueOf(t.ProtocolFactory)})
		if len(result) != 1 {
			return nil, errors.New("func must return 1 result")
		}
		return &ServiceWrapper{service: result[0].Interface(), transport: transport}, err
	}

}
