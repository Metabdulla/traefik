package auth

import (
	"context"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"github.com/containous/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

const (
	appKeyTypeName = "AppKeyAuth"
	appKeyHeader   = "appKey" //Todo
)

type appKeyAuth struct {
	next         http.Handler
	headerField  string
	removeHeader bool
	name         string
	appKeys      map[string][]string
}

// NewBasic creates a basicAuth middleware.
func NewAppKeyAuth(ctx context.Context, next http.Handler, authConfig dynamic.AppKeyAuth, name string) (http.Handler, error) {
	middlewares.GetLogger(ctx, name, basicTypeName).Debug("Creating middleware")
	auth := &appKeyAuth{
		next:         next,
		appKeys:      make(map[string][]string),
		headerField:  authConfig.HeaderField,
		removeHeader: authConfig.RemoveHeader,
		name:         name,
	}

	return auth, nil
}

func (b *appKeyAuth) GetTracingInformation() (string, ext.SpanKindEnum) {
	return b.name, tracing.SpanKindNoneEnum
}

func (a *appKeyAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := middlewares.GetLogger(req.Context(), a.name, appKeyTypeName)
	appKey := req.Header.Get(appKeyHeader)
	if ok := a.verify(appKey,"todo"); ok {
		a.next.ServeHTTP(rw, req)
	} else {
		logger.Debug("Authentication failed")
		tracing.SetErrorWithEvent(req, "Authentication failed")
		a.RequireAuth(rw, req)
	}
}

func (a *appKeyAuth) RequireAuth(w http.ResponseWriter, r *http.Request) {
	//todo fix content type
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("not authorized"))
}

func (a *appKeyAuth) verify(appKey string, apiName string ) bool {
	//todo
	if appKey =="" {
		return false
	}

	 keys := a.appKeys[apiName]
	 if len(keys) <1 {
	 	return false
	 }
	 for _, key := range keys {
	 	if key ==appKey {
	 		return true
		}
	}

	return true
}
