package nacos

import (
	"fmt"
	"github.com/containous/traefik/v2/pkg/config/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/sirupsen/logrus"
	"net/url"
)

var nacosClient *NacosClient

type NacosClient struct {
	iNamingClient naming_client.INamingClient
	iConfigClient  config_client.IConfigClient
}

//NacosClientInit  init nacos client when start up
func  NacosClientInit (config nacos.Nacos) *NacosClient{
	if !config.Enable{
		return nil
	}
	var serverConfigs = config.ServerConfig()
	if len(serverConfigs)==0 {
		return nil
	}
	clientConfig:= config.ClientConfig()

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		//"serverConfigs": []constant.ServerConfig{
		//	{
		//		IpAddr: "172.28.152.105",
		//		Port:   32148,
		//	},
		//},
		"serverConfigs":serverConfigs,
		"clientConfig": *clientConfig,
	})
	if namingClient ==nil {
		logrus.WithError(err).Error("client init failed")
		return nil
	}
	nacosClient = &NacosClient{}
	nacosClient.iNamingClient = namingClient
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		//"serverConfigs": []constant.ServerConfig{
		//	{
		//		IpAddr: "172.28.152.105",
		//		Port:   32148,
		//	},
		//},
		"serverConfigs":serverConfigs,
		"clientConfig": *clientConfig,
	})
	if configClient ==nil {
		logrus.WithError(err).Error("client init failed")
		return nil
	}
	nacosClient.iConfigClient = configClient
	return nacosClient
}




func (n *NacosClient)INamingClient (nacos nacos.Nacos)(naming_client.INamingClient) {
   return n.iNamingClient
}

func (n *NacosClient)IConfigClient (nacos nacos.Nacos)(config_client.IConfigClient) {
	return n.iConfigClient
}

func (n*NacosClient)SelectOneHealthyInstance(serviceName  string ,clusterName  []string )( url string ,err error) {
	instance ,err := n.iNamingClient.SelectOneHealthyInstance( vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		Clusters:  clusterName ,
	})
	if err!=nil {
		return "",err
	}
	return fmt.Sprintf("http://%s:%d",instance.Ip,instance.Port),nil
}



func (n*NacosClient)CheckService (serviceName  string ,clusterName  []string ,url *url.URL)  bool {
	urls,_ := n.GetServers(serviceName,clusterName)
	if len(urls) ==0 {
		return false
	}

	for _,u:= range urls {
		if sameURL(u, url) {
			return true
		}
	}
	return false
}


func (n*NacosClient)GetServers(serviceName  string ,clusterName  []string) ([]*url.URL, error){
	instances ,err := n.iNamingClient.SelectInstances( vo.SelectInstancesParam{ Clusters:clusterName,
		ServiceName:serviceName,
		HealthyOnly:true,
	})
	if err!=nil {
		logrus.WithError(err).Debug("get service  failed")
		return nil,err
	}
	var urls []*url.URL
	for _, v:= range instances {
		url ,err := url.Parse(fmt.Sprintf("http://%s:%d", v.Ip, v.Port))
		if err!=nil {
			logrus.WithError(err).Errorf("parser url  failed")
			continue
		}
		urls = append(urls, url)

	}
	return urls,nil
}

func sameURL(a, b *url.URL) bool {
	return a.Path == b.Path && a.Host == b.Host && a.Scheme == b.Scheme
}