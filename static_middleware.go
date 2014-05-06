package traffic

import (
	"net/http"
	"os"
	"path/filepath"
)

type StaticMiddleware struct {
	publicPath string
}

func (middleware *StaticMiddleware) ServeHTTP(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	path := filepath.Join(middleware.publicPath, r.URL.Path)
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		w.Header().Del("Content-Type")
		http.ServeFile(w, r.Request, path)

		return
	}

	if nextMiddleware := next(); nextMiddleware != nil {
		nextMiddleware.ServeHTTP(w, r, next)
	}
}

func NewStaticMiddleware(publicPath string) *StaticMiddleware {
	middleware := &StaticMiddleware{
		publicPath: publicPath,
	}

	return middleware
}
