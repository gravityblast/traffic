package traffic

import (
  "net/http"
)

type HttpMethod string

type Router struct {
  routes map[HttpMethod][]*Route
}

func (router *Router) Add(method HttpMethod, route *Route) {
  router.routes[method] = append(router.routes[method], route)
}

func (router *Router) Get(path string, handler HttpHandleFunc) {
  route := NewRoute(path, handler)
  router.Add(HttpMethod("GET"), route)
}

func (router *Router) ServeHTTP(r http.ResponseWriter, req *http.Request) {
  for _, route := range router.routes[HttpMethod(req.Method)] {
    values, ok := route.Match(req.URL.Path)
    if ok {
      newValues := req.URL.Query()
      for k, v := range values {
        newValues[k] = v
      }

      req.URL.RawQuery = newValues.Encode()

      route.Handler(r, req)
      return
    }
  }
}

func New() *Router {
  router := &Router{}
  router.routes = make(map[HttpMethod][]*Route)
  return router
}

