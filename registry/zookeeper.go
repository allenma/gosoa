//zookeeper.go
package registry

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZKConfig struct {
	Servers        []string
	SessionTimeout time.Duration
	RootDir        string
}

func NewZKRegistry(servers []string, timeout time.Duration) (registry *ZKRegistry, err error) {
	zkConfig := &ZKConfig{
		Servers:        servers,
		SessionTimeout: timeout,
		RootDir:        "/soaregistry",
	}
	conn, _, err := zk.Connect(servers, timeout)
	createZKNodeIfNotExist(conn, zkConfig.RootDir, 0, nil)
	return &ZKRegistry{Conn: conn, ZKConf: zkConfig}, err
}

type ZKRegistry struct {
	ZKConf *ZKConfig
	Conn   *zk.Conn
}

func (z *ZKRegistry) ListServices() (services []string, err error) {
	services, _, err = z.Conn.Children(z.ZKConf.RootDir)
	return
}

func (z *ZKRegistry) ListServiceProviders(serviceName string) (providers []ProviderInfo, err error) {
	providers, err = z.doGetProviders(z.getZkPath(serviceName))
	return
}

func (z *ZKRegistry) UpdateServiceProvider(serviceName string, provider ProviderInfo) (err error) {
	jsonReq, err := json.Marshal(provider)
	zkpath := z.getZkPath(serviceName) + "/" + provider.Addr
	if err == nil {
		_, err = z.Conn.Set(zkpath, []byte(jsonReq), -1)
	}

	return
}

func (z *ZKRegistry) Register(serviceName string, providerInfo ProviderInfo) (err error) {
	serviceZKPath := fmt.Sprintf("%s/%s", z.ZKConf.RootDir, serviceName)
	zkpath := fmt.Sprintf("%s/%s", serviceZKPath, providerInfo.Addr)
	jsonReq, err := json.Marshal(providerInfo)
	if err != nil {
		return
	}

	err = createZKNodeIfNotExist(z.Conn, serviceZKPath, 0, nil)
	err = createZKNodeIfNotExist(z.Conn, zkpath, zk.FlagEphemeral, []byte(jsonReq))
	return
}

func createZKNodeIfNotExist(conn *zk.Conn, zkPath string, flags int32, data []byte) (err error) {
	exist, _, err := conn.Exists(zkPath)
	if !exist {
		conn.Create(zkPath, data, flags, zk.WorldACL(zk.PermAll))
	}
	return
}

func (z *ZKRegistry) Discover(serviceName string, callback RegChangeCallback) (providers []ProviderInfo, err error) {
	zkpath := z.getZkPath(serviceName)
	providers, err = z.doGetProviders(zkpath)
	go func() {
		for {
			addresses, _, eventChan, err := z.Conn.ChildrenW(zkpath)
			fmt.Println("children and watch:", addresses)
			if err != nil {
				fmt.Println("error when get children")
				return
			}
			providers := make([]ProviderInfo, 0, len(addresses))
			aggDataChan := make(chan zk.Event)
			for _, addr := range addresses {
				providerPath := zkpath + "/" + addr
				_, _, dataEventChan, err := z.Conn.GetW(providerPath)
				if err == nil {
					go func() {
						evt := <-dataEventChan
						aggDataChan <- evt
					}()
				}
			}
			select {
			case <-eventChan:
			case <-aggDataChan:
				providers, err = z.doGetProviders(serviceName)
				if err == nil {
					callback(providers, err)
				}

			}

		}
	}()
	return
}

func (z *ZKRegistry) getZkPath(serviceName string) string {
	return fmt.Sprintf("%s/%s", z.ZKConf.RootDir, serviceName)
}

func (z *ZKRegistry) doGetProviders(zkpath string) (providers []ProviderInfo, err error) {
	addresses, _, err := z.Conn.Children(zkpath)
	if err == nil {
		providers := make([]ProviderInfo, 0, len(addresses))
		for _, addr := range addresses {
			providerPath := zkpath + "/" + addr
			providerData, _, err := z.Conn.Get(providerPath)
			providerInfo := ProviderInfo{}
			err = json.Unmarshal(providerData, &providerInfo)
			if err != nil {
				fmt.Println("error when unmashal provider:", addr, " error:", err.Error()) // TODO change to use log
			} else {
				providers = append(providers, providerInfo)
			}
		}
		return providers, err
	}
	return
}
