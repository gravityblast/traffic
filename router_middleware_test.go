package traffic

import (
	"fmt"
	assert "github.com/pilu/miniassert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeNotFoundHandlerContainer struct {
	callsCount int
}

func (f *fakeNotFoundHandlerContainer) Handler(w ResponseWriter, r *Request) {
	f.callsCount++
}

func newTestRequest(method, path string) (ResponseWriter, *httptest.ResponseRecorder, *Request) {
	r, _ := http.NewRequest(method, path, nil)
	request := newRequest(r)
	recorder := httptest.NewRecorder()

	env := make(map[string]interface{})
	responseWriter := newResponseWriter(recorder, &env)

	return responseWriter, recorder, request
}

func newTestRouterMiddleware() *RouterMiddleware {
	router := &Router{}
	router.routes = make(map[HttpMethod][]*Route)
	router.beforeFilters = make([]HttpHandleFunc, 0)
	router.middlewares = make([]Middleware, 0)
	routerMiddleware := &RouterMiddleware{router}

	return routerMiddleware
}

func TestRouterMiddleware_NotFound(t *testing.T) {
	routerMiddleware := newTestRouterMiddleware()
	fakeNotFound := new(fakeNotFoundHandlerContainer)
	routerMiddleware.router.NotFoundHandler = fakeNotFound.Handler

	responseWriter, recorder, request := newTestRequest("GET", "/")
	routerMiddleware.ServeHTTP(responseWriter, request, func() Middleware { return nil })

	// checks that the router middleware calls router.handleNotFound
	assert.Equal(t, 1, fakeNotFound.callsCount)
	assert.Equal(t, "", string(recorder.Body.Bytes()))
}

func TestRouterMiddleware_Found(t *testing.T) {
	routerMiddleware := newTestRouterMiddleware()
	responseWriter, recorder, request := newTestRequest("GET", "/")

	testRootHandler := func(w ResponseWriter, r *Request) {
		fmt.Fprint(w, "Hello World")
	}

	routerMiddleware.router.Get("/", testRootHandler)

	routerMiddleware.ServeHTTP(responseWriter, request, func() Middleware { return nil })
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "Hello World", string(recorder.Body.Bytes()))
}
