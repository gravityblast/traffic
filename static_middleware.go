package traffic

import (
	"net/http"
)

type StaticMiddleware struct {
	publicPath string
}

func (middleware *StaticMiddleware) ServeHTTP(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	callNext := func() {
		if nextMiddleware := next(); nextMiddleware != nil {
			nextMiddleware.ServeHTTP(w, r, next)
		}
	}

	dir := http.Dir(middleware.publicPath)
	path := r.URL.Path

	file, err := dir.Open(path)
	if err != nil {
		callNext()
		return
	}
	defer file.Close()

	if info, err := file.Stat(); err == nil && !info.IsDir() {
		w.Header().Del("Content-Type")
		http.ServeContent(w, r.Request, path, info.ModTime(), file)
		return
	}

	callNext()
}

func NewStaticMiddleware(publicPath string) *StaticMiddleware {
	middleware := &StaticMiddleware{
		publicPath: publicPath,
	}

	return middleware
}
