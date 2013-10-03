package traffic

import (
  "os"
  "net/http"
  "path/filepath"
)

type StaticMiddleware struct {
  publicPath string
}

func (middleware *StaticMiddleware) ServeHTTP(w ResponseWriter, r *http.Request, next NextMiddlewareFunc) (ResponseWriter, *http.Request) {
  path := filepath.Join(middleware.publicPath, r.URL.Path)
  if info, err := os.Stat(path); err == nil && !info.IsDir() {
    w.Header().Del("Content-Type")
    http.ServeFile(w, r, path)

    return w, r
  }

  if nextMiddleware := next(); nextMiddleware != nil {
    w, r = nextMiddleware.ServeHTTP(w, r, next)
  }

  return w, r
}

func NewStaticMiddleware(publicPath string) *StaticMiddleware {
  middleware := &StaticMiddleware{
    publicPath: publicPath,
  }

  return middleware
}

