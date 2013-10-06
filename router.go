package traffic

import (
  "net/http"
  "os"
  "log"
  "fmt"
  "github.com/pilu/config"
)

type HttpMethod string

type BeforeFilterFunc func(ResponseWriter, *http.Request) bool

type NextMiddlewareFunc func() Middleware

type Middleware interface {
  ServeHTTP(ResponseWriter, *http.Request, NextMiddlewareFunc) (ResponseWriter, *http.Request)
}

type Router struct {
  routes map[HttpMethod][]*Route
  NotFoundHandler HttpHandleFunc
  beforeFilters []BeforeFilterFunc
  middlewares []Middleware
  env map[string]interface{}
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

func (router *Router) handleNotFound (w ResponseWriter, r *http.Request) {
  if router.NotFoundHandler != nil {
    router.NotFoundHandler(w, r)
  } else {
    fmt.Fprint(w, "404 page not found")
  }
}

func (router *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
  w := newAppResponseWriter(rw, &router.env)
  w.Header().Set("Content-Type", "text/html")

  nextMiddlewareFunc := router.MiddlewareEnumerator()
  if nextMiddleware := nextMiddlewareFunc(); nextMiddleware != nil {
    nextMiddleware.ServeHTTP(w, r, nextMiddlewareFunc)
  }

  if w.StatusCode() == http.StatusNotFound {
    router.handleNotFound(w, r)
  }
}

func (router *Router) AddMiddleware(middleware Middleware) {
  router.middlewares = append([]Middleware{middleware}, router.middlewares...)
}

func (router *Router) SetVar(key string, value interface{}) {
  router.env[key] = value
}

func (router *Router) GetVar(key string) interface{} {
  value := router.env[key]
  if value != nil {
    return value
  }

  return GetVar(key)
}

func (router *Router) loadConfigurationsFromFile(path, env string) {
  mainSectionName := "main"
  sections, err := config.ParseFile(path, mainSectionName)
  if err != nil {
    panic(err)
  }

  for section, options := range sections {
    if section == mainSectionName || section == env {
      for key, value := range options {
        router.SetVar(key, value)
      }
    }
  }
}

func New() *Router {
  router := &Router{}
  router.routes = make(map[HttpMethod][]*Route)
  router.beforeFilters = make([]BeforeFilterFunc, 0)
  router.middlewares = make([]Middleware, 0)
  router.env = make(map[string]interface{})

  routerMiddleware := &RouterMiddleware{ router }
  router.AddMiddleware(routerMiddleware)

  // Logger
  logger := log.New(os.Stderr, "", log.LstdFlags)
  router.SetVar("logger", logger)

  // Environment
  env := Env()

  // Add useful middlewares for development
  if env == EnvDevelopment {
    // Static middleware
    router.AddMiddleware(NewStaticMiddleware(PublicPath()))

    // Logger middleware
    loggerMiddleware := &LoggerMiddleware{
      router: router,
      logger: logger,
    }
    router.AddMiddleware(loggerMiddleware)

    // ShowErrors middleware
    router.AddMiddleware(&ShowErrorsMiddleware{})
  }

  // configuration
  configFile := ConfigFilePath()
  if _, err := os.Stat(configFile); err == nil {
    router.loadConfigurationsFromFile(configFile, env)
  }

  initTemplateManager()

  return router
}

