
package traffic

import (
  "testing"
  assert "github.com/pilu/miniassert"
)

func TestPathSegmentToRegexpSegment(t *testing.T) {
  tests := [][]string{
    {
      "foo",
      "foo",
    },
    {
      ":foo",
      "(?P<foo>[^/#?]+)",
    },
  }

  for _, pair := range tests {
    pathSegment := pair[0]
    expectedRegexpSegment := pair[1]
    assert.Equal(t, expectedRegexpSegment, pathSegmentToRegexpSegment(pathSegment))
  }
}

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
  }

  for _, pair := range tests {
    pathSegment := pair[0]
    expectedRegexpSegment := pair[1]
    assert.Equal(t, expectedRegexpSegment, pathToRegexpString(pathSegment))
  }
}
