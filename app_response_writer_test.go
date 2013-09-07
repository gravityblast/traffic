package traffic

import (
  "testing"
  "net/http"
  "net/http/httptest"
  assert "github.com/pilu/miniassert"
)

func newTestAppResponseWriter(globalEnv *map[string]interface{}) *AppResponseWriter {
  recorder := httptest.NewRecorder()
  arw := newAppResponseWriter(
    recorder,
    globalEnv,
  )
  return arw
}

func TestAppResponseWriter(t *testing.T) {
  globalEnv := make(map[string]interface{})
  arw := newTestAppResponseWriter(&globalEnv)
  assert.Equal(t, http.StatusOK, arw.statusCode)
  assert.Equal(t, 0, len(arw.env))
  assert.Equal(t, &globalEnv, arw.globalEnv)
}

func TestAppResponseWriter_SetVar(t *testing.T) {
  globalEnv := make(map[string]interface{})
  arw := newTestAppResponseWriter(&globalEnv)
  arw.SetVar("foo", "bar")
  assert.Equal(t, "bar", arw.env["foo"])
}

func TestAppResponseWriter_GetVar(t *testing.T) {
  globalEnv := map[string]interface{} {
    "global_foo": "global_bar",
  }
  arw := newTestAppResponseWriter(&globalEnv)
  arw.env["foo"] = "bar"

  assert.Equal(t, "bar", arw.GetVar("foo"))
  assert.Equal(t, "global_bar", arw.GetVar("global_foo"))

  arw.SetVar("global_foo", "local_bar")
  assert.Equal(t, "global_bar", (*arw.globalEnv)["global_foo"])
  assert.Equal(t, "local_bar", arw.GetVar("global_foo"))
}
