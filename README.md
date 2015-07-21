## SOA framework for go
go实现的简单soa框架。支持thrift rpc框架，支持服务的注册与自动发现，目前支持redis和zookeeper作为服务注册中心，可以配置客户端负载均衡方式，目前支持基于加权的轮询和随机两种负载均衡方式，可以设置客户端自动重试次数。提供了简单的管理restful api，可以查看服务有哪些提供者，mark down/up 服务提供者，调整服务提供者的权重等等。

## 使用
### thrift 代码生成
下载thrift: https://thrift.apache.org/download  
下载thrift go client: go get git.apache.org/thrift.git/lib/go/thrift   
生成thrift go 客户端代码：   
thrift -r --gen go:package_prefix=github.com/allenma/gosoa/sample/ tutorial.thrift

具体可参考wiki: [https://thrift.apache.org/tutorial/go]

### server端

	handler := NewCalculatorHandler()
    processor := tutorial.NewCalculatorProcessor(handler)
	reg := registry.NewRedisRegistry("localhost:6379")
	exporter := export.NewThriftExporter("localhost:9090", reg)
	err := exporter.Export("calculator",processor)
	if err!=nil {
		fmt.Println("error when export service",err.Error())
		return
	}
	
	
上面代码是使用redis作为服务注册中心，要用zookeeper只需更改一行代码即可
	
	reg := registry.NewZKRegistry([]string{"localhost:2181"}, 5*time.Second)
	

### client端

    reg := registry.NewRedisRegistry("localhost:6379")
	//	reg,_ := registry.NewZKRegistry([]string{"192.168.148.128:2181"},10*time.Second)
	tclient := client.NetThriftClient(reg, "calculator", tutorial.NewCalculatorClientFactory)
	var a int32 = 5
	var b int32 = 6
	result, err := tclient.Execute("Add", a, b)

	if err == nil {
		fmt.Println("  add result=", result.(int32))
	} 
	_,err = tclient.Execute("Ping")
	

具体代码可以参考包： github.com/allenma/gosoa/sample/, 其中server目录是service provider导出的代码，client目录是service consumer代码   


### 管理

- 更改github.com/allenma/gosoa/manage/manage.go中服务注册中心的地址(后续会放在配置中)
- 运行管理程序:go run $GOPATH/github.com/allenma/gosoa/manage/manage.go  
 
获取所有服务：http://localhost:8080/service  

获取服务的所有提供者：http://localhost:8080/service/${SERVICE_NAME}/providers  

更改服务提供者信息：curl -d "{\"addr\":\"localhost:9091\",\"status\":1,\"weight\":5}" http://localhost:8080/service/${SERVICE_NAME}/updateprovider  


## TODO
- 客户端熔断机制
- 完善服务管理
- 服务监控
- 支持别的rpc框架，比如grpc,etc.