package traffic

import (
  "github.com/pilu/fresh/runner/runnerutils"
)

type buildErrorsMiddleware struct {}

func (middleware buildErrorsMiddleware) ServeHTTP(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
  if runnerutils.HasErrors() {
    runnerutils.RenderError(w)
    return
  }

  if nextMiddleware := next(); nextMiddleware != nil {
    nextMiddleware.ServeHTTP(w, r, next)
  }
}

func newBuildErrorsMiddleware() *buildErrorsMiddleware {
  return &buildErrorsMiddleware{}
}
