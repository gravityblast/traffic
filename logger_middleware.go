package traffic

import (
  "net/http"
  "log"
  "time"
)

type LoggerMiddleware struct {
  router *Router
  logger *log.Logger
}

func (loggerMiddleware *LoggerMiddleware) ServeHTTP(w ResponseWriter, r *http.Request, next NextMiddlewareFunc) (ResponseWriter, *http.Request) {
  startTime := time.Now()

  if nextMiddleware := next(); nextMiddleware != nil {
    w, r = nextMiddleware.ServeHTTP(w, r, next)
  }

  duration := time.Since(startTime).Seconds()

  loggerMiddleware.logger.Printf(`[%.6f] %d "%s"`, duration, w.StatusCode(), r.URL.Path)

  return w, r
}
