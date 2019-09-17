package dynamic

//appKey for every api

type AppKeyAuth struct {
	AppKeys        AppKeys  `json:"appKeys,omitempty" toml:"appKeys,omitempty" yaml:"appKeys,omitempty"`
	RemoveHeader bool   `json:"removeHeader,omitempty" toml:"removeHeader,omitempty" yaml:"removeHeader,omitempty"`
	HeaderField  string `json:"headerField,omitempty" toml:"headerField,omitempty" yaml:"headerField,omitempty" export:"true"`
}

// Users holds a list of users
type AppKeys []string


type NacosService struct {
	Clusters []string    `description:"nacos  services  options" json:"clusters,omitempty" toml:"clusters,omitempty" yaml:"clusters,omitempty"`
	ServiceName  string `description:"nacos  services  options" json:"serviceName,omitempty" toml:"serviceName,omitempty" yaml:"serviceName,omitempty"`
	Group  string  `description:"nacos  services  options" json:"group,omitempty" toml:"group,omitempty" yaml:"group,omitempty"`
}