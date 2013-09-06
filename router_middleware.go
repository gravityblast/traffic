package traffic

import (
  "net/http"
)

type RouterMiddleware struct {
  router *Router
}

func (routerMiddleware *RouterMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next NextMiddlewareFunc) (http.ResponseWriter, *http.Request) {
  for _, route := range routerMiddleware.router.routes[HttpMethod(r.Method)] {
    values, ok := route.Match(r.URL.Path)
    if ok {
      newValues := r.URL.Query()
      for k, v := range values {
        newValues[k] = v
      }

      r.URL.RawQuery = newValues.Encode()

      continueAfterBeforeFilter := true

      filters := append(routerMiddleware.router.beforeFilters, route.beforeFilters...)

      for _, beforeFilter := range filters {
        continueAfterBeforeFilter = beforeFilter(w, r)
        if !continueAfterBeforeFilter {
          break
        }
      }

      if continueAfterBeforeFilter {
        route.Handler(w, r)
      }

      return w, r
    }
  }

  w.WriteHeader(http.StatusNotFound)

  return w, r
}
