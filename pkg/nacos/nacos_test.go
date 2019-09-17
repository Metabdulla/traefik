package nacos

import (
	"fmt"
	"github.com/containous/traefik/v2/pkg/config/nacos"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"
	"time"
)

func TestGetService(t *testing.T) {
	serverConfig := nacos.NacosServerConfig{
		IpAddr: "172.28.152.102",
		Port:   32148,
	}
	var config  nacos.Nacos
	config.Enable = true
	config.Servers = append(config.Servers,serverConfig)
	fmt.Println(config)
	nacosClient:= NacosClientInit(config)
	if nacosClient==nil {
		t.Fatal("init fail")
	}
	ok ,err := nacosClient.iNamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "172.28.152.225",
		Port:        8295,
		ServiceName: "demoHaha.go",
		Weight:      5,
		ClusterName: "b",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if err!=nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("register service fail")
	}


	//example.ExampleServiceClient_GetService(client)
	service, _ := nacosClient.iNamingClient.GetService(vo.GetServiceParam{
		ServiceName: "demoHaha.go",
		Clusters:    []string{"b"},
	})
	fmt.Println(utils.ToJsonString(service))
	instances ,_:=  nacosClient.iNamingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "demoHaha.go",
		Clusters:    []string{"b"},
	})
	fmt.Println(utils.ToJsonString(instances))
	instance ,_:=  nacosClient.iNamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "demoHaha.go",
		Clusters:    []string{"b"},
	})
	fmt.Println(utils.ToJsonString(instance))
	result ,_:= nacosClient.SelectOneHealthyInstance("demoHaha.go",[]string{"b"})
	fmt.Println(result)
	for {
		time.Sleep(10*time.Second)
	}
}
