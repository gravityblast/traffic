package traffic

import (
  "testing"
  "regexp"
  "net/http"
  "net/url"
  "reflect"
  assert "github.com/pilu/miniassert"
)

func httpHandlerExample(r ResponseWriter, req *http.Request) {}

func TestNewRoute(t *testing.T) {
  path := "/categories/:category_id/posts/:id"
  route := NewRoute(path, httpHandlerExample)
  assert.Type(t, "*traffic.Route", route)
  assert.Equal(t, path, route.Path)

  expectedPathRegexp := regexp.MustCompile("^/categories/(?P<category_id>[^/#?]+)/posts/(?P<id>[^/#?]+)$")
  assert.Equal(t, expectedPathRegexp, route.PathRegexp)
}

func TestRoute_Match(t *testing.T) {
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
    {
      "/pages/:page_path*",
      "/pages/foo/bar/baz",
      "page_path=foo%2Fbar%2Fbaz",
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

func TestRoute_Match_WithOptionalSegments(t *testing.T) {
  routePath := "((/sites/:site_id)?/categories/:category_id)?/posts/:id"
  tests := [][]string {
    {
      "/sites/foo/categories/bar/posts/baz",
      "category_id=bar&id=baz&site_id=foo",
    },
    {
      "/categories/bar/posts/baz",
      "category_id=bar&id=baz",
    },
    {
      "/posts/baz",
      "id=baz",
    },
  }
  route := NewRoute(routePath, httpHandlerExample)
  for _, opts := range tests {
    requestPath   := opts[0]
    expectedQuery := opts[1]
    values, ok := route.Match(requestPath)
    assert.True(t, ok)
    assert.Equal(t, expectedQuery, values.Encode())
  }
}

func TestRoute_AddBeforeFilterToRoute(t *testing.T) {
  route := NewRoute("/", httpHandlerExample)
  assert.Equal(t, 0, len(route.beforeFilters))
  filterA := BeforeFilterFunc(func(w ResponseWriter, r *http.Request) bool { return true })
  filterB := BeforeFilterFunc(func(w ResponseWriter, r *http.Request) bool { return true })

  route.AddBeforeFilter(filterA)
  assert.Equal(t, 1, len(route.beforeFilters))
  route.AddBeforeFilter(filterB)
  assert.Equal(t, 2, len(route.beforeFilters))

  assert.Equal(t, reflect.ValueOf(filterA), reflect.ValueOf(route.beforeFilters[0]))
  assert.Equal(t, reflect.ValueOf(filterB), reflect.ValueOf(route.beforeFilters[1]))
}
