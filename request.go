package traffic

import (
	"net/http"
	"net/url"
)

type Request struct {
	*http.Request
}

func (r Request) Params() url.Values {
	return r.Request.URL.Query()
}

func (r Request) Param(key string) string {
	return r.Params().Get(key)
}

func newRequest(r *http.Request) *Request {
	return &Request{
		Request: r,
	}
}
