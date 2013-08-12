package traffic

import (
  "net/http"
)

type HttpMethod string

type BeforeFilterFunc func(http.ResponseWriter, *http.Request)

type Router struct {
  routes map[HttpMethod][]*Route
  NotFoundHandler HttpHandleFunc
  BeforeFilter BeforeFilterFunc
}

func (router *Router) Add(method HttpMethod, path string, handler HttpHandleFunc) {
  route := NewRoute(path, handler)
  router.routes[method] = append(router.routes[method], route)
}

func (router *Router) Get(path string, handler HttpHandleFunc) {
  router.Add(HttpMethod("GET"), path, handler)
  router.Add(HttpMethod("HEAD"), path, handler)
}

func (router *Router) Post(path string, handler HttpHandleFunc) {
  router.Add(HttpMethod("POST"), path, handler)
}

func (router *Router) Delete(path string, handler HttpHandleFunc) {
  router.Add(HttpMethod("DELETE"), path, handler)
}

func (router *Router) Put(path string, handler HttpHandleFunc) {
  router.Add(HttpMethod("PUT"), path, handler)
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

      if router.BeforeFilter != nil {
        router.BeforeFilter(w, r)
      }

      route.Handler(w, r)
      return
    }
  }

  if router.NotFoundHandler != nil {
    router.NotFoundHandler(w, r)
  } else {
    http.NotFound(w, r)
  }
}

func New() *Router {
  router := &Router{}
  router.routes = make(map[HttpMethod][]*Route)
  return router
}

