// redis
package registry

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisConfig struct {
	Address          string
	ConnectTimeoutMs int
	WriteTimeoutMs   int
	ReadTimeoutMs    int
	MaxIdle          int
	MaxActive        int
	IdleTimeoutS     int
	Password         string
}

type redisRegistry struct {
	pool *redis.Pool
}

const REDIS_NS = "soa:"

func NewRedisRegistryWithConfig(conf *RedisConfig) (registry Registry) {
	pool := &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeoutS) * time.Second,
		MaxActive:   conf.MaxActive,
		Dial: func() (redis.Conn, error) {
			var c redis.Conn
			var err error
			c, err = redis.DialTimeout("tcp", conf.Address,
				time.Duration(conf.ConnectTimeoutMs)*time.Millisecond,
				time.Duration(conf.ReadTimeoutMs)*time.Millisecond,
				time.Duration(conf.WriteTimeoutMs)*time.Millisecond)
			if err != nil {
				return nil, err
			}
			//password authentication
			if len(conf.Password) > 0 {
				if _, err_pass := c.Do("AUTH", conf.Password); err_pass != nil {
					c.Close()
				}
			}
			return c, err
		},
	}
	return &redisRegistry{pool}
}

func NewRedisRegistry(addr string) (registry Registry) {
	conf := &RedisConfig{
		Address:          addr,
		ConnectTimeoutMs: 500,
		WriteTimeoutMs:   500,
		ReadTimeoutMs:    500,
		MaxIdle:          5,
		MaxActive:        10,
		IdleTimeoutS:     1800,
	}
	return NewRedisRegistryWithConfig(conf)
}

func (r *redisRegistry) ListServices() (services []string, err error) {
	conn := r.pool.Get()
	defer conn.Close()

	servicesArr, err := redis.Strings(conn.Do("KEYS", REDIS_NS+"*"))
	if err == nil {
		services = make([]string, 0, len(servicesArr))
		for _, serviceName := range servicesArr {
			services = append(services, strings.TrimLeft(serviceName, REDIS_NS))
		}
	}
	return
}

func (r *redisRegistry) ListServiceProviders(serviceName string) (providers []ProviderInfo, err error) {
	providers, err = r.doFindAll(serviceName)
	return
}

func (r *redisRegistry) UpdateServiceProvider(serviceName string, provider ProviderInfo) (err error) {
	conn := r.pool.Get()
	defer conn.Close()

	jsonReq, err := json.Marshal(provider)
	if err == nil {
		_, err = conn.Do("HSET", r.getFullServiceName(serviceName), provider.Addr, jsonReq)
	}

	return
}

func (r *redisRegistry) Register(serviceName string, providerInfo ProviderInfo) (err error) {
	conn := r.pool.Get()
	defer conn.Close()

	var input_params []interface{}
	jsonReq, err := json.Marshal(providerInfo)
	if err != nil {
		return
	}

	input_params = append(input_params, r.getFullServiceName(serviceName))
	input_params = append(input_params, providerInfo.Addr)
	input_params = append(input_params, string(jsonReq))
	_, err = redis.Int64(conn.Do("HSET", input_params...))
	return
}

func (r *redisRegistry) Discover(serviceName string, callback RegChangeCallback) (infos []ProviderInfo, err error) {

	infos, err = r.doFindAll(serviceName)

	// update address cache every minute
	timer := time.NewTicker(1 * time.Minute)
	go func() {
		for range timer.C {
			newVal, err := r.doFindAll(serviceName)
			callback(newVal, err)
		}
	}()
	return infos, err
}

func (r *redisRegistry) doFindAll(serviceName string) (infos []ProviderInfo, err error) {
	conn := r.pool.Get()
	defer conn.Close()

	jsonStrs, err := redis.Strings(conn.Do("HVALS", r.getFullServiceName(serviceName)))
	if err != nil {
		return
	}
	providers := make([]ProviderInfo, 0, len(jsonStrs))
	for _, jsonStr := range jsonStrs {
		p := ProviderInfo{}
		json.Unmarshal([]byte(jsonStr), &p)
		providers = append(providers, p)
	}

	return providers, nil
}

func (r *redisRegistry) getFullServiceName(serviceName string) string {
	return REDIS_NS + serviceName
}
