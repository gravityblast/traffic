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
  router.Add(HttpMethod("HEAD"), route)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  for _, route := range router.routes[HttpMethod(r.Method)] {
    values, ok := route.Match(r.URL.Path)
    if ok {
      newValues := r.URL.Query()
      for k, v := range values {
        newValues[k] = v
      }

      r.URL.RawQuery = newValues.Encode()

      route.Handler(w, r)
      return
    }
  }

  http.NotFound(w, r)
}

func New() *Router {
  router := &Router{}
  router.routes = make(map[HttpMethod][]*Route)
  return router
}

