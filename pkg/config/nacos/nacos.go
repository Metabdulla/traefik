package nacos

import 	(
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

type Nacos struct {
	Servers []NacosServerConfig `description:"nacos servers configuration options" json:"servers,omitempty" toml:"servers,omitempty" yaml:"servers,omitempty" export:"true"`
	Client *NacosClientConfig `description:"nacos client configuration options" json:"client,omitempty" toml:"client,omitempty" yaml:"client,omitempty" export:"true"`
	Enable  bool `description:"nacos  configuration options" json:"enable,omitempty" toml:"enable,omitempty" yaml:"enable,omitempty" export:"true"`
}

type NacosServerConfig struct {
	ContextPath string `description:"ContextPath  configuration options" json:"contextPath,omitempty" toml:"contextPath,omitempty" yaml:"contextPath,omitempty" export:"true"`
	IpAddr      string  `description:"IpAddr  configuration options" json:"ipAddr,omitempty" toml:"ipAddr,omitempty" yaml:"ipAddr,omitempty" export:"true"`
	Port        uint64 `description:"Port  configuration options" json:"port,omitempty" toml:"port,omitempty" yaml:"port,omitempty" export:"true"`
}

type NacosClientConfig struct {
	TimeoutMs            uint64 `description:"TimeoutMs  configuration options" json:"timeoutMs,omitempty" toml:"timeoutMs,omitempty" yaml:"timeoutMs,omitempty" export:"true"`
	ListenInterval       uint64 `description:"ListenInterval  configuration options" json:"listenInterval,omitempty" toml:"listenInterval,omitempty" yaml:"listenInterval,omitempty" export:"true"`
	BeatInterval         int64 `description:"BeatInterval  configuration options" json:"beatInterval,omitempty" toml:"beatInterval,omitempty" yaml:"beatInterval,omitempty" export:"true"`
	NamespaceId          string `description:"NamespaceId  configuration options" json:"namespaceId,omitempty" toml:"namespaceId,omitempty" yaml:"namespaceId,omitempty" export:"true"`
	Endpoint             string `description:"Endpoint  configuration options" json:"endpoint,omitempty" toml:"endpoint,omitempty" yaml:"endpoint,omitempty" export:"true"`
	AccessKey            string `description:"AccessKey  configuration options" json:"accessKey,omitempty" toml:"accessKey,omitempty" yaml:"accessKey,omitempty" export:"true"`
	SecretKey            string `description:"SecretKey  configuration options" json:"secretKey,omitempty" toml:"secretKey,omitempty" yaml:"secretKey,omitempty" export:"true"`
	CacheDir             string `description:"CacheDir  configuration options" json:"cacheDir,omitempty" toml:"cacheDir,omitempty" yaml:"cacheDir,omitempty" export:"true"`
	LogDir               string `description:"LogDir  configuration options" json:"logDir,omitempty" toml:"logDir,omitempty" yaml:"logDir,omitempty" export:"true"`
	UpdateThreadNum      int `description:"UpdateThreadNum  configuration options" json:"updateThreadNum,omitempty" toml:"updateThreadNum,omitempty" yaml:"updateThreadNum,omitempty" export:"true"`
	NotLoadCacheAtStart  bool `description:"NotLoadCacheAtStart  configuration options" json:"notLoadCacheAtStart,omitempty" toml:"notLoadCacheAtStart,omitempty" yaml:"notLoadCacheAtStart,omitempty" export:"true"`
	UpdateCacheWhenEmpty bool `description:"UpdateCacheWhenEmpty  configuration options" json:"updateCacheWhenEmpty,omitempty" toml:"updateCacheWhenEmpty,omitempty" yaml:"updateCacheWhenEmpty,omitempty" export:"true"`
	OpenKMS              bool `description:"OpenKMS  configuration options" json:"openKMS,omitempty" toml:"openKMS,omitempty" yaml:"openKMS,omitempty" export:"true"`
	RegionId             string `description:"RegionId  configuration options" json:"regionId,omitempty" toml:"regionId,omitempty" yaml:"regionId,omitempty" export:"true"`
}

func (c*Nacos)ServerConfig() ([]constant.ServerConfig){
	if c==nil {
		return nil
	}
	if len(c.Servers) == 0 {
		return nil
	}
	var serverConfigs []constant.ServerConfig
	for _,v := range c.Servers {
		config:= constant.ServerConfig{
			ContextPath: v.ContextPath,
			IpAddr:      v.IpAddr,
			Port:        v.Port,
		}
		serverConfigs = append(serverConfigs,config)
	}
		return serverConfigs
}

func (n*Nacos)ClientConfig() (*constant.ClientConfig){
	if n==nil {
		return nil
	}
	if n.Client ==nil {
		n.Client = NewDefaultClientConfig()
	}
	c:= n.Client
	return &constant.ClientConfig{
		TimeoutMs            :c.TimeoutMs,
		ListenInterval       :c.ListenInterval,
		BeatInterval         :c.BeatInterval,
		NamespaceId          :c.NamespaceId,
		Endpoint             :c.Endpoint,
		AccessKey            :c.AccessKey,
		SecretKey            :c.SecretKey,
		CacheDir             :c.CacheDir,
		LogDir               :c.LogDir,
		UpdateThreadNum      :c.UpdateThreadNum,
		NotLoadCacheAtStart  :c.NotLoadCacheAtStart,
		UpdateCacheWhenEmpty :c.UpdateCacheWhenEmpty,
		OpenKMS              :c.OpenKMS,
		RegionId             :c.RegionId,
	}
}

func NewDefaultClientConfig() *NacosClientConfig {
	return &NacosClientConfig{
		TimeoutMs:      10 * 1000, //http请求超时时间，单位毫秒
		ListenInterval: 30 * 1000, //监听间隔时间，单位毫秒（仅在ConfigClient中有效）
		BeatInterval:   5 * 1000, //心跳间隔时间，单位毫秒（仅在ServiceClient中有效）
		NamespaceId:       "public", //nacos命名空间
		Endpoint:          "" ,//获取nacos节点ip的服务地址
		CacheDir:         "nacos_cache", //缓存目录
		LogDir:         "nacos_log", //日志目录
		UpdateThreadNum:   20, //更新服务的线程数
		NotLoadCacheAtStart: true, //在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: true, //当服务列表为空时是否更新本地缓存，true--更新,false--不更新
	}
}