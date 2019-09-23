package cryptSign

//RequestBodyRaw  make sure that the filed is sorted
type RequestBodyRaw struct {
	AppKey  string `json:"appKey"`
	BizContent  string `json:"bizContent"`
	ServiceName string `json:"serviceName"`
	Timestamp   string `json:"timestamp"`
	Version string `json:"version"`
}

type RequestBody struct {
	RequestBodyRaw
	SignValue  string `json:"signValue"`

}

