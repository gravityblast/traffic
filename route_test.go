package traffic

import (
  "testing"
  "regexp"
  "net/http"
  "net/url"
  assert "github.com/pilu/miniassert"
)

func httpHandlerExample(r http.ResponseWriter, req *http.Request) {}

func TestNewRoute(t *testing.T) {
  path := "/categories/:category_id/posts/:id"
  route := NewRoute(path, httpHandlerExample)
  assert.Type(t, "*traffic.Route", route)
  assert.Equal(t, path, route.Path)

  expectedPathRegexp := regexp.MustCompile("^/categories/(?P<category_id>[^/#?]+)/posts/(?P<id>[^/#?]+)$")
  assert.Equal(t, expectedPathRegexp, route.PathRegexp)
}

func TestMatch(t *testing.T) {
  tests := [][]string {
    {
      "/",
      "/",
      "",
    },
    {
      "/:foo/?",
      "/bar",
      "foo=bar",
    },
    {
      "/:foo/?",
      "/bar/", // with trailing slash
      "foo=bar",
    },
    {
      "/categories/:category_id/posts/:id",
      "/categories/foo/posts/bar",
      "category_id=foo&id=bar",
    },
  }

  for _, opts := range tests {
    routePath     := opts[0]
    requestPath   := opts[1]
    expectedQuery := opts[2]

    route := NewRoute(routePath, httpHandlerExample)
    values, ok := route.Match(requestPath)
    assert.True(t, ok)
    assert.Equal(t, expectedQuery, values.Encode())
  }

  route := NewRoute("/foo", httpHandlerExample)
  values, ok := route.Match("/bar")
  assert.False(t, ok)
  assert.Equal(t, values, make(url.Values))
}
