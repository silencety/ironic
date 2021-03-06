package local

import (
	"api/server/httputils"
	dkrouter "api/server/router"
)

// It implements router.Route.
type localRoute struct {
	method  string
	path    string
	handler httputils.APIFunc
}

type router struct{
	routers []dkrouter.Route
}

// Handler returns the APIFunc to let the server wrap it in middlewares
func (l localRoute) Handler() httputils.APIFunc {
	return l.handler
}

// Method returns the http method that the route responds to.
func (l localRoute) Method() string {
	return l.method
}

// Path returns the subpath where the route responds to.
func (l localRoute) Path() string {
	return l.path
}

// NewRoute initializes a new local router for the reouter
func NewRoute(method, path string, handler httputils.APIFunc) dkrouter.Route {
	return localRoute{method, path, handler}
}

// NewGetRoute initializes a new route with the http method GET.
func NewGetRoute(path string, handler httputils.APIFunc) dkrouter.Route {
	return NewRoute("GET", path, handler)
}

// NewPostRoute initializes a new route with the http method POST.
func NewPostRoute(path string, handler httputils.APIFunc) dkrouter.Route {
	return NewRoute("POST", path, handler)
}

// NewPutRoute initializes a new route with the http method PUT.
func NewPutRoute(path string, handler httputils.APIFunc) dkrouter.Route {
	return NewRoute("PUT", path, handler)
}

// NewDeleteRoute initializes a new route with the http method DELETE.
func NewDeleteRoute(path string, handler httputils.APIFunc) dkrouter.Route {
	return NewRoute("DELETE", path, handler)
}

// NewOptionsRoute initializes a new route with the http method OPTIONS
func NewOptionsRoute(path string, handler httputils.APIFunc) dkrouter.Route {
	return NewRoute("OPTIONS", path, handler)
}

// NewHeadRoute initializes a new route with the http method HEAD.
func NewHeadRoute(path string, handler httputils.APIFunc) dkrouter.Route {
	return NewRoute("HEAD", path, handler)
}

