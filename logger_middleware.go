package traffic

import (
	"time"
)

type LoggerMiddleware struct {
	router *Router
}

func (loggerMiddleware *LoggerMiddleware) ServeHTTP(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	startTime := time.Now()

	if nextMiddleware := next(); nextMiddleware != nil {
		nextMiddleware.ServeHTTP(w, r, next)
	}

	duration := time.Since(startTime).Seconds()

	Logger().Printf(`[%.6f] %d "%s"`, duration, w.StatusCode(), r.URL.Path)
}
