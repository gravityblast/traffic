package traffic

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/pilu/config"
)

type HttpMethod string

type ErrorHandlerFunc func(ResponseWriter, *Request, interface{})

type NextMiddlewareFunc func() Middleware

type HttpHandleFunc func(ResponseWriter, *Request)

type Middleware interface {
	ServeHTTP(ResponseWriter, *Request, NextMiddlewareFunc)
}

type Router struct {
	NotFoundHandler HttpHandleFunc
	ErrorHandler    ErrorHandlerFunc
	routes          map[HttpMethod][]*Route
	beforeFilters   []HttpHandleFunc
	middlewares     []Middleware
	env             map[string]interface{}
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

func (router *Router) Add(method HttpMethod, path string, handlers ...HttpHandleFunc) *Route {
	route := NewRoute(path, handlers...)
	router.addRoute(method, route)

	return route
}

func (router *Router) addRoute(method HttpMethod, route *Route) {
	router.routes[method] = append(router.routes[method], route)
}

func (router *Router) Get(path string, handlers ...HttpHandleFunc) *Route {
	route := router.Add(HttpMethod("GET"), path, handlers...)
	router.addRoute(HttpMethod("HEAD"), route)

	return route
}

func (router *Router) Post(path string, handlers ...HttpHandleFunc) *Route {
	return router.Add(HttpMethod("POST"), path, handlers...)
}

func (router *Router) Delete(path string, handlers ...HttpHandleFunc) *Route {
	return router.Add(HttpMethod("DELETE"), path, handlers...)
}

func (router *Router) Put(path string, handlers ...HttpHandleFunc) *Route {
	return router.Add(HttpMethod("PUT"), path, handlers...)
}

func (router *Router) Patch(path string, handlers ...HttpHandleFunc) *Route {
	return router.Add(HttpMethod("PATCH"), path, handlers...)
}

func (router *Router) AddBeforeFilter(beforeFilters ...HttpHandleFunc) *Router {
	router.beforeFilters = append(router.beforeFilters, beforeFilters...)

	return router
}

func (router *Router) handleNotFound(w ResponseWriter, r *Request) {
	if router.NotFoundHandler != nil {
		router.NotFoundHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 page not found")
	}
}

func (router *Router) handlePanic(w ResponseWriter, r *Request, err interface{}) {
	if router.ErrorHandler != nil {
		w.WriteHeader(http.StatusInternalServerError)
		router.ErrorHandler(w, r, err)
	} else {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	const size = 4096
	stack := make([]byte, size)
	stack = stack[:runtime.Stack(stack, false)]

	logger.Printf("%v\n", err)
	logger.Printf("%s\n", string(stack))
}

func (router *Router) ServeHTTP(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	w := newResponseWriter(httpResponseWriter, &router.env)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	r := newRequest(httpRequest)

	defer func() {
		if recovered := recover(); recovered != nil {
			router.handlePanic(w, r, recovered)
		}
	}()

	nextMiddlewareFunc := router.MiddlewareEnumerator()
	if nextMiddleware := nextMiddlewareFunc(); nextMiddleware != nil {
		nextMiddleware.ServeHTTP(w, r, nextMiddlewareFunc)
	}
}

func (router *Router) Use(middleware Middleware) {
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

func addDevelopmentMiddlewares(router *Router) {
	// Static middleware
	router.Use(NewStaticMiddleware(PublicPath()))

	// Logger middleware
	loggerMiddleware := &LoggerMiddleware{
		router: router,
	}
	router.Use(loggerMiddleware)

	// ShowErrors middleware
	router.Use(&ShowErrorsMiddleware{})
}

func (router *Router) Run() {
	address := fmt.Sprintf("%s:%d", Host(), Port())
	Logger().Printf("Starting in %s on %s", Env(), address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal(err)
	}
}

func (router *Router) RunSSL(certFile, keyFile string) {
	address := fmt.Sprintf("%s:%d", Host(), Port())
	Logger().Printf("Starting in %s on %s", Env(), address)
	err := http.ListenAndServeTLS(address, certFile, keyFile, router)
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfigurationsFromFile(path, env string) {
	mainSectionName := "main"
	sections, err := config.ParseFile(path, mainSectionName)
	if err != nil {
		panic(err)
	}

	for section, options := range sections {
		if section == mainSectionName || section == env {
			for key, value := range options {
				SetVar(key, value)
			}
		}
	}
}

func init() {
	env = make(map[string]interface{})
	SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	// configuration
	configFile := ConfigFilePath()
	if _, err := os.Stat(configFile); err == nil {
		loadConfigurationsFromFile(configFile, Env())
	}
}

func New() *Router {
	router := &Router{
		routes:        make(map[HttpMethod][]*Route),
		beforeFilters: make([]HttpHandleFunc, 0),
		middlewares:   make([]Middleware, 0),
		env:           make(map[string]interface{}),
	}

	routerMiddleware := &RouterMiddleware{router}
	router.Use(routerMiddleware)

	// Environment
	env := Env()

	// Add useful middlewares for development
	if env == EnvDevelopment {
		addDevelopmentMiddlewares(router)
	}

	initTemplateManager()

	return router
}
