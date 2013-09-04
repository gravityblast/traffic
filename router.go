package traffic

import (
  "net/http"
  "os"
  "log"
  "fmt"
)

type HttpMethod string

type BeforeFilterFunc func(http.ResponseWriter, *http.Request) bool

type Middleware interface {
  ServeHTTP(http.ResponseWriter, *http.Request, func() Middleware) (http.ResponseWriter, *http.Request)
}

type Router struct {
  routes map[HttpMethod][]*Route
  NotFoundHandler HttpHandleFunc
  beforeFilters []BeforeFilterFunc
  middlewares []Middleware
}

func (router Router) MiddlewareEnumerator() func() Middleware {
  index := 0
  next := func() Middleware {
    if len(router.middlewares) > index {
      nextMiddleware := router.middlewares[index]
      index++
      return nextMiddleware
    }

    return nil
  }

  return next
}

func (router *Router) Add(method HttpMethod, path string, handler HttpHandleFunc) *Route {
  route := NewRoute(path, handler)
  router.addRoute(method, route)

  return route
}

func (router *Router) addRoute(method HttpMethod, route *Route) {
  router.routes[method] = append(router.routes[method], route)
}

func (router *Router) Get(path string, handler HttpHandleFunc) *Route {
  route := router.Add(HttpMethod("GET"), path, handler)
  router.addRoute(HttpMethod("HEAD"), route)

  return route
}

func (router *Router) Post(path string, handler HttpHandleFunc) *Route {
  return router.Add(HttpMethod("POST"), path, handler)
}

func (router *Router) Delete(path string, handler HttpHandleFunc) *Route {
  return router.Add(HttpMethod("DELETE"), path, handler)
}

func (router *Router) Put(path string, handler HttpHandleFunc) *Route {
  return router.Add(HttpMethod("PUT"), path, handler)
}

func (router *Router) Patch(path string, handler HttpHandleFunc) *Route {
  return router.Add(HttpMethod("PATCH"), path, handler)
}

func (router *Router) AddBeforeFilter(beforeFilter BeforeFilterFunc) *Router {
  router.beforeFilters = append(router.beforeFilters, beforeFilter)

  return router
}

func (router *Router) handleNotFound (w http.ResponseWriter, r *http.Request) {
  if router.NotFoundHandler != nil {
    router.NotFoundHandler(w, r)
  } else {
    fmt.Fprint(w, "404 page not found")
  }
}

func (router *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
  w := &AppResponseWriter{
    ResponseWriter: rw,
  }

  nextMiddlewareFunc := router.MiddlewareEnumerator()
  if nextMiddleware := nextMiddlewareFunc(); nextMiddleware != nil {
    nextMiddleware.ServeHTTP(w, r, nextMiddlewareFunc)
  }

  if w.StatusCode() == http.StatusNotFound {
    router.handleNotFound(w, r)
  }
}

func New() *Router {
  router := &Router{}
  router.routes = make(map[HttpMethod][]*Route)
  router.beforeFilters = make([]BeforeFilterFunc, 0)
  router.middlewares = make([]Middleware, 0)

  loggerMiddleware := &LoggerMiddleware{
    router: router,
    logger: log.New(os.Stderr, "", log.LstdFlags),
  }
  router.middlewares = append(router.middlewares, loggerMiddleware)

  routerMiddleware := &RouterMiddleware{ router }
  router.middlewares = append(router.middlewares, routerMiddleware)

  return router
}

