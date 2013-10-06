package traffic

import (
  "time"
  "net/http"
)

type LoggerMiddleware struct {
  router *Router
}

func (loggerMiddleware *LoggerMiddleware) ServeHTTP(w ResponseWriter, r *http.Request, next NextMiddlewareFunc) (ResponseWriter, *http.Request) {
  startTime := time.Now()

  if nextMiddleware := next(); nextMiddleware != nil {
    w, r = nextMiddleware.ServeHTTP(w, r, next)
  }

  duration := time.Since(startTime).Seconds()

  Logger().Printf(`[%.6f] %d "%s"`, duration, w.StatusCode(), r.URL.Path)

  return w, r
}
