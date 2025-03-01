package axon

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	*httprouter.Router
}

func NewRouter() *Router {
	return &Router{
		Router: httprouter.New(),
	}
}

func (r *Router) GET(path string, handle httprouter.Handle) {
	r.Handle("GET", path, NoCacheHandle(handle))
}

func (r *Router) HEAD(path string, handle httprouter.Handle) {
	r.Handle("HEAD", path, NoCacheHandle(handle))
}

func (r *Router) POST(path string, handle httprouter.Handle) {
	r.Handle("POST", path, NoCacheHandle(handle))
}

func (r *Router) PUT(path string, handle httprouter.Handle) {
	r.Handle("PUT", path, NoCacheHandle(handle))
}

func (r *Router) PATCH(path string, handle httprouter.Handle) {
	r.Handle("PATCH", path, NoCacheHandle(handle))
}

func (r *Router) DELETE(path string, handle httprouter.Handle) {
	r.Handle("DELETE", path, NoCacheHandle(handle))
}

func (r *Router) Handle(method, path string, handle httprouter.Handle) {
	r.Router.Handle(method, path, NoCacheHandle(handle))
}

func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Router.HandlerFunc(method, path, NoCacheHandlerFunc(handler))
}

func (r *Router) Handler(method, path string, handler http.Handler) {
	r.Router.Handler(method, path, NoCacheHandler(handler))
}
