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

func (loggerMiddleware *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, nextMiddlewareFunc func() Middleware) (http.ResponseWriter, *http.Request) {
  nextMiddleware := nextMiddlewareFunc()
  startTime := time.Now()

  if nextMiddleware != nil {
    w, r = nextMiddleware.ServeHTTP(w, r, nextMiddlewareFunc)
  }

  duration := time.Since(startTime).Seconds()

  if arw, ok := w.(*AppResponseWriter); ok {
    loggerMiddleware.logger.Printf(`[%.6f] %d "%s"`, duration, arw.StatusCode(), r.URL.Path)
  }

  return w, r
}
