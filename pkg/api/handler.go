package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/middlewares/auth"
	"net/http"
	"strconv"
	"strings"

	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/version"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

const (
	defaultPerPage = 100
	defaultPage    = 1
)

const nextPageHeader = "X-Next-Page"

type serviceInfoRepresentation struct {
	*runtime.ServiceInfo
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
}

// RunTimeRepresentation is the configuration information exposed by the API handler.
type RunTimeRepresentation struct {
	Routers     map[string]*runtime.RouterInfo        `json:"routers,omitempty"`
	Middlewares map[string]*runtime.MiddlewareInfo    `json:"middlewares,omitempty"`
	Services    map[string]*serviceInfoRepresentation `json:"services,omitempty"`
	TCPRouters  map[string]*runtime.TCPRouterInfo     `json:"tcpRouters,omitempty"`
	TCPServices map[string]*runtime.TCPServiceInfo    `json:"tcpServices,omitempty"`
}

type pageInfo struct {
	startIndex int
	endIndex   int
	nextPage   int
}

// Handler serves the configuration and status of Traefik on API endpoints.
type Handler struct {
	dashboard bool
	debug     bool
	// runtimeConfiguration is the data set used to create all the data representations exposed by the API.
	runtimeConfiguration *runtime.Configuration
	staticConfig         static.Configuration
	// statistics           *types.Statistics
	// stats                *thoasstats.Stats // FIXME stats
	// StatsRecorder         *middlewares.StatsRecorder // FIXME stats
	dashboardAssets *assetfs.AssetFS

	BasicAuth  *dynamic.BasicAuth
}

// New returns a Handler defined by staticConfig, and if provided, by runtimeConfig.
// It finishes populating the information provided in the runtimeConfig.
func New(staticConfig static.Configuration, runtimeConfig *runtime.Configuration) *Handler {
	rConfig := runtimeConfig
	if rConfig == nil {
		rConfig = &runtime.Configuration{}
	}

	return &Handler{
		dashboard: staticConfig.API.Dashboard,
		// statistics:           staticConfig.API.Statistics,
		dashboardAssets:      staticConfig.API.DashboardAssets,
		runtimeConfiguration: rConfig,
		staticConfig:         staticConfig,
		debug:                staticConfig.API.Debug,
		BasicAuth: staticConfig.API.BasicAuth,
	}
}

// Append add api routes on a router
func (h Handler) Append(router *mux.Router) {
	if h.debug {
		DebugHandler{}.Append(router)
	}
	if h.BasicAuth!=nil {
		ctx := context.Background()
		getRuntimeConfiguration ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getRuntimeConfiguration),*h.BasicAuth,"getRuntimeConfiguration")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/rawdata").HandlerFunc(getRuntimeConfiguration.ServeHTTP)

		// Experimental endpoint
		getOverview ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getOverview),*h.BasicAuth,"getOverview")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/overview").HandlerFunc(getOverview.ServeHTTP)

		getEntryPoints ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getEntryPoints),*h.BasicAuth,"getEntryPoints")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/entrypoints").HandlerFunc(getEntryPoints.ServeHTTP)

		getEntryPoint ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getEntryPoint),*h.BasicAuth,"getEntryPoint")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/entrypoints/{entryPointID}").HandlerFunc(getEntryPoint.ServeHTTP)

		getRouters ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getRouters),*h.BasicAuth,"getRouters")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/http/routers").HandlerFunc(getRouters.ServeHTTP)
		getRouter ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getRouter),*h.BasicAuth,"getRouter")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/http/routers/{routerID}").HandlerFunc(getRouter.ServeHTTP)

		getServices ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getServices),*h.BasicAuth,"getServices")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/http/services").HandlerFunc(getServices.ServeHTTP)

		getService ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getService),*h.BasicAuth,"getService")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/http/services/{serviceID}").HandlerFunc(getService.ServeHTTP)

		getMiddlewares ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getMiddlewares),*h.BasicAuth,"getMiddlewares")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/http/middlewares").HandlerFunc(getMiddlewares.ServeHTTP)

		getMiddleware ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getMiddleware),*h.BasicAuth,"getMiddleware")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/http/middlewares/{middlewareID}").HandlerFunc(getMiddleware.ServeHTTP)

		getTCPRouters ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getTCPRouters),*h.BasicAuth,"getTCPRouters")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/tcp/routers").HandlerFunc(getTCPRouters.ServeHTTP)

		getTCPRouter ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getTCPRouter),*h.BasicAuth,"getTCPRouter")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/tcp/routers/{routerID}").HandlerFunc(getTCPRouter.ServeHTTP)
		getTCPServices ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getTCPServices),*h.BasicAuth,"getTCPServices")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/tcp/services").HandlerFunc(getTCPServices.ServeHTTP)
		getTCPService ,err := auth.NewBasic(ctx,http.HandlerFunc(h.getTCPService),*h.BasicAuth,"getTCPService")
		if err!=nil {
			panic(err)
		}
		router.Methods(http.MethodGet).Path("/api/tcp/services/{serviceID}").HandlerFunc(getTCPService.ServeHTTP)

	}else {

		router.Methods(http.MethodGet).Path("/api/rawdata").HandlerFunc(h.getRuntimeConfiguration)

		// Experimental endpoint
		router.Methods(http.MethodGet).Path("/api/overview").HandlerFunc(h.getOverview)

		router.Methods(http.MethodGet).Path("/api/entrypoints").HandlerFunc(h.getEntryPoints)
		router.Methods(http.MethodGet).Path("/api/entrypoints/{entryPointID}").HandlerFunc(h.getEntryPoint)

		router.Methods(http.MethodGet).Path("/api/http/routers").HandlerFunc(h.getRouters)
		router.Methods(http.MethodGet).Path("/api/http/routers/{routerID}").HandlerFunc(h.getRouter)
		router.Methods(http.MethodGet).Path("/api/http/services").HandlerFunc(h.getServices)
		router.Methods(http.MethodGet).Path("/api/http/services/{serviceID}").HandlerFunc(h.getService)
		router.Methods(http.MethodGet).Path("/api/http/middlewares").HandlerFunc(h.getMiddlewares)
		router.Methods(http.MethodGet).Path("/api/http/middlewares/{middlewareID}").HandlerFunc(h.getMiddleware)

		router.Methods(http.MethodGet).Path("/api/tcp/routers").HandlerFunc(h.getTCPRouters)
		router.Methods(http.MethodGet).Path("/api/tcp/routers/{routerID}").HandlerFunc(h.getTCPRouter)
		router.Methods(http.MethodGet).Path("/api/tcp/services").HandlerFunc(h.getTCPServices)
		router.Methods(http.MethodGet).Path("/api/tcp/services/{serviceID}").HandlerFunc(h.getTCPService)
	}

	// FIXME stats
	// health route
	// router.Methods(http.MethodGet).Path("/health").HandlerFunc(p.getHealthHandler)

	version.Handler{}.Append(router)

	if h.dashboard {
		DashboardHandler{Assets: h.dashboardAssets}.Append(router)
	}
}

func (h Handler) getRuntimeConfiguration(rw http.ResponseWriter, request *http.Request) {

	siRepr := make(map[string]*serviceInfoRepresentation, len(h.runtimeConfiguration.Services))
	for k, v := range h.runtimeConfiguration.Services {
		siRepr[k] = &serviceInfoRepresentation{
			ServiceInfo:  v,
			ServerStatus: v.GetAllStatus(),
		}
	}

	result := RunTimeRepresentation{
		Routers:     h.runtimeConfiguration.Routers,
		Middlewares: h.runtimeConfiguration.Middlewares,
		Services:    siRepr,
		TCPRouters:  h.runtimeConfiguration.TCPRouters,
		TCPServices: h.runtimeConfiguration.TCPServices,
	}

	rw.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func pagination(request *http.Request, max int) (pageInfo, error) {
	perPage, err := getIntParam(request, "per_page", defaultPerPage)
	if err != nil {
		return pageInfo{}, err
	}

	page, err := getIntParam(request, "page", defaultPage)
	if err != nil {
		return pageInfo{}, err
	}

	startIndex := (page - 1) * perPage
	if startIndex != 0 && startIndex >= max {
		return pageInfo{}, fmt.Errorf("invalid request: page: %d, per_page: %d", page, perPage)
	}

	endIndex := startIndex + perPage
	if endIndex >= max {
		endIndex = max
	}

	nextPage := 1
	if page*perPage < max {
		nextPage = page + 1
	}

	return pageInfo{startIndex: startIndex, endIndex: endIndex, nextPage: nextPage}, nil
}

func getIntParam(request *http.Request, key string, defaultValue int) (int, error) {
	raw := request.URL.Query().Get(key)
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("invalid request: %s: %d", key, value)
	}
	return value, nil
}

func getProviderName(id string) string {
	return strings.SplitN(id, "@", 2)[1]
}
