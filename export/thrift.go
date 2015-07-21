// thrift.go
package export

import (
	"crypto/tls"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/allenma/gosoa/registry"
)

type ThriftExporter struct {
	Provider registry.ProviderInfo
	Reg      registry.Registry
	Config   *ThriftConfig
}

type ThriftConfig struct {
	TransFactory    thrift.TTransportFactory
	ProtocolFactory thrift.TProtocolFactory
	Secure          bool
	CertFile        string
	KeyFile         string
}

func NewThriftExporter(addr string, reg registry.Registry) *ThriftExporter {
	return &ThriftExporter{
		Provider: registry.ProviderInfo{Addr: addr, Status: 0, Weight: 5},
		Reg:      reg,
		Config: &ThriftConfig{
			TransFactory:    thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()),
			ProtocolFactory: thrift.NewTBinaryProtocolFactoryDefault(),
			Secure:          false,
		},
	}
}

func (t *ThriftExporter) Export(serviceName string, processor thrift.TProcessor) (err error) {
	var transport thrift.TServerTransport
	if t.Config.Secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair(t.Config.CertFile, t.Config.KeyFile); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(t.Provider.Addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(t.Provider.Addr)
	}

	if err != nil {
		return err
	}
	server := thrift.NewTSimpleServer4(processor, transport, t.Config.TransFactory, t.Config.ProtocolFactory)

	err = t.Reg.Register(serviceName, t.Provider)
	if err != nil {
		fmt.Println("error when register service", err.Error())
		return
	}
	fmt.Println("Starting the simple server... on ", t.Provider.Addr)
	return server.Serve()
}
