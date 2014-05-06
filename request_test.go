package traffic

import (
	assert "github.com/pilu/miniassert"
	"net/http"
	"net/url"
	"testing"
)

func TestRequest_Params(t *testing.T) {
	r, _ := http.NewRequest("GET", "http:///example.com?f=foo&b=bar", nil)
	request := newRequest(r)
	params := make(url.Values)
	params.Set("f", "foo")
	params.Set("b", "bar")

	assert.Equal(t, params, request.Params())
}

func TestRequest_Param(t *testing.T) {
	r, _ := http.NewRequest("GET", "http:///example.com?f=foo&b=bar", nil)
	request := newRequest(r)

	assert.Equal(t, "foo", request.Param("f"))
	assert.Equal(t, "bar", request.Param("b"))
}
