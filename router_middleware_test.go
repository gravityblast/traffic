package traffic

import (
  "fmt"
  "testing"
  "net/http"
  "net/http/httptest"
  assert "github.com/pilu/miniassert"
)

func newTestRequest(method, path string) (ResponseWriter, *httptest.ResponseRecorder, *http.Request) {
  request, _ := http.NewRequest(method, path, nil)
  recorder := httptest.NewRecorder()

  env := make(map[string]interface{})
  appResponseWriter := newAppResponseWriter(recorder, &env)

  return appResponseWriter, recorder, request
}

func newTestRouterMiddleware() *RouterMiddleware {
  router := &Router{}
  router.routes = make(map[HttpMethod][]*Route)
  router.beforeFilters = make([]HttpHandleFunc, 0)
  router.middlewares = make([]Middleware, 0)
  routerMiddleware := &RouterMiddleware{ router }

  return routerMiddleware
}

func TestRouterMiddleware_NotFound(t *testing.T) {
  routerMiddleware := newTestRouterMiddleware()
  responseWriter, recorder, request := newTestRequest("GET", "/")
  routerMiddleware.ServeHTTP(responseWriter, request, func() Middleware { return nil })

  assert.Equal(t, 1, 1)
  assert.Equal(t, 404, recorder.Code)
  assert.Equal(t, "", string(recorder.Body.Bytes()))
}

func TestRouterMiddleware_Found(t *testing.T) {
  routerMiddleware := newTestRouterMiddleware()
  responseWriter, recorder, request := newTestRequest("GET", "/")

  testRootHandler := func (w ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello World")
  }

  routerMiddleware.router.Get("/", testRootHandler)

  routerMiddleware.ServeHTTP(responseWriter, request, func() Middleware { return nil })
  assert.Equal(t, 200, recorder.Code)
  assert.Equal(t, "Hello World", string(recorder.Body.Bytes()))
}
