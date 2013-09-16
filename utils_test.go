
package traffic

import (
  "testing"
  "os"
  assert "github.com/pilu/miniassert"
)

func resetGlobalEnv() {
  env = make(map[string]interface{})
}

func TestSetVar(t *testing.T) {
  resetGlobalEnv()
  SetVar("foo", "bar")
  assert.Equal(t, "bar", env["foo"])
  resetGlobalEnv()
}

func TestGetVar(t *testing.T) {
  resetGlobalEnv()
  env["foo-2"] = "bar-2"
  assert.Equal(t, "bar-2", GetVar("foo-2"))

  assert.Nil(t, GetVar("os_foo"))
  os.Setenv("TRAFFIC_OS_FOO", "bar")
  assert.Equal(t, "bar", GetVar("os_foo"))


  resetGlobalEnv()
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
