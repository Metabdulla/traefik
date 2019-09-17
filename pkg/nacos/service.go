package nacos

import (
	"context"
	"errors"
	"fmt"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"github.com/containous/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/roundrobin"
	"github.com/vulcand/oxy/utils"
	"net/http"
	"net/url"
)


const (
	typeName = "NacosService"
)

// ipWhiteLister is a middleware that provides Checks of the Requesting IP against a set of Whitelists
type NacosService struct {
	next        http.Handler
	name string
	serviceName string
	cluster     []string
	group        string
	stickySession        *StickySession
	nacosClient    *NacosClient

}

// New builds a new NacosService given a serviceName and cluster
func New(ctx context.Context, next http.Handler, config dynamic.NacosService, name string ,session *StickySession) (*NacosService, error) {
	logger := middlewares.GetLogger(ctx, name, typeName)
	logger.Debug("Creating middleware")

	if config.ServiceName == "" || len(config.Clusters) == 0 {
		return nil, errors.New("ServiceName or Clusters  is empty, NacosService not created")
	}

	logger.Debugf("Setting up NacosService with name: %s", config.ServiceName)
	return &NacosService{
		next:        next,
		serviceName:   config.ServiceName,
		cluster:config.Clusters,
		group: config.Group,
		name:name,
		stickySession:session,
		nacosClient:nacosClient,
	}, nil
}

func (wl *NacosService) GetTracingInformation() (string, ext.SpanKindEnum) {
	return wl.name, tracing.SpanKindNoneEnum
}

func (n *NacosService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := middlewares.GetLogger(req.Context(), n.name, typeName)
	logEntry := log.WithField("Request", utils.DumpHttpRequest(req))
	logEntry.Debug("nacos service : begin ServeHttp on request")
	defer logEntry.Debug("nacos service:  completed ServeHttp on request")


	// make shallow copy of request before chaning anything to avoid side effects
	newReq := *req
	stuck := false
	if n.stickySession != nil {
		cookieURL, present, err := n.stickySession.GetBackend(&newReq, n.Servers())
		if cookieURL!= nil{
			present  = n.nacosClient.CheckService(n.serviceName,n.cluster,cookieURL)

		}

		if err != nil {
			log.Warnf("nacos service: error using server from cookie: %v", err)
		}

		if present {
			newReq.URL = cookieURL
			stuck = true
		}
	}

	if !stuck {
		url, err := n.NextServer()
		if err != nil {
			logMessage := fmt.Sprintf("service address  not found request %+v: %v %v %v", req, err,n.serviceName,n.cluster)
			log.Debug(logMessage)
			tracing.SetErrorWithEvent(req, logMessage)
			//todo
			notFound(log, w)
			return

		}

		if n.stickySession != nil {
			n.stickySession.StickBackend(url, &w)
		}
		newReq.URL = url
	}


		// log which backend URL we're sending this request to
		log.WithFields(logrus.Fields{"Request": utils.DumpHttpRequest(req), "ForwardURL": newReq.URL}).Debugf("nacos service: Forwarding this request to URL")

	// Emit event to a listener if one exists

	n.next.ServeHTTP(w, &newReq)
}

func notFound(logger logrus.FieldLogger, rw http.ResponseWriter) {
	statusCode := http.StatusNotFound

	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(http.StatusText(statusCode)))
	if err != nil {
		logger.Error(err)
	}
}

func (n *NacosService) Servers() []*url.URL{
	urls,_ := n.nacosClient.GetServers(n.serviceName,n.cluster)
	//todo
	return urls
}

func (n *NacosService) RemoveServer(u *url.URL) error {
	//todo  unregister instance ??
 	return nil
}

func (n*NacosService) UpsertServer(u *url.URL, options ...roundrobin.ServerOption) error{
	////todo  register instance ??
 	return nil
}

func (n *NacosService)NextServer()(*url.URL,error){
	urlStr,err:=  n.nacosClient.SelectOneHealthyInstance(n.serviceName,n.cluster)
	if err!=nil {
		return nil,err
	}
	return  url.Parse(urlStr)
}

