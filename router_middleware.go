package traffic

import (
	"net/http"
)

type RouterMiddleware struct {
	router *Router
}

func (routerMiddleware *RouterMiddleware) ServeHTTP(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	for _, route := range routerMiddleware.router.routes[HttpMethod(r.Method)] {
		values, ok := route.Match(r.URL.Path)
		if ok {
			newValues := r.URL.Query()
			for k, v := range values {
				newValues[k] = v
			}

			r.URL.RawQuery = newValues.Encode()

			handlers := append(routerMiddleware.router.beforeFilters, route.beforeFilters...)
			handlers = append(handlers, route.Handlers...)

			for _, handler := range handlers {
				handler(w, r)
				if w.Written() {
					break
				}
			}

			if w.StatusCode() == http.StatusNotFound && !w.BodyWritten() {
				routerMiddleware.router.handleNotFound(w, r)
			}

			return
		}
	}

	routerMiddleware.router.handleNotFound(w, r)
}
