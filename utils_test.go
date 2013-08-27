
package traffic

import (
  "testing"
  assert "github.com/pilu/miniassert"
)

func TestPathToRegexpString(t *testing.T) {
  tests := [][]string{
    {
      "/",
      "^/$",
    },
    {
      "/foo/bar",
      "^/foo/bar$",
    },
    {
      "/foo/bar",
      "^/foo/bar$",
    },
    {
      "/:foo/bar/:baz",
      "^/(?P<foo>[^/#?]+)/bar/(?P<baz>[^/#?]+)$",
    },
    {
      "(/categories/:category_id)?/posts/:id",
      "^(/categories/(?P<category_id>[^/#?]+))?/posts/(?P<id>[^/#?]+)$",
    },
  }

  for _, pair := range tests {
    pathSegment := pair[0]
    expectedRegexpSegment := pair[1]
    assert.Equal(t, expectedRegexpSegment, pathToRegexpString(pathSegment))
  }
}
