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
  routerEnv := make(map[string]interface{})
  arw := newTestAppResponseWriter(&routerEnv)
  assert.Equal(t, http.StatusOK, arw.statusCode)
  assert.Equal(t, 0, len(arw.env))
  assert.Equal(t, &routerEnv, arw.routerEnv)
}

func TestAppResponseWriter_SetVar(t *testing.T) {
  globalEnv := make(map[string]interface{})
  arw := newTestAppResponseWriter(&globalEnv)
  arw.SetVar("foo", "bar")
  assert.Equal(t, "bar", arw.env["foo"])
}

func TestAppResponseWriter_GetVar(t *testing.T) {
  resetGlobalEnv()

  routerEnv := map[string]interface{} {}
  arw := newTestAppResponseWriter(&routerEnv)

  env["global-foo"] = "global-bar"
  assert.Equal(t, "global-bar", arw.GetVar("global-foo"))

  routerEnv["global-foo"] = "router-bar"
  assert.Equal(t, "router-bar", arw.GetVar("global-foo"))

  arw.env["global-foo"] = "local-bar"
  assert.Equal(t, "local-bar", arw.GetVar("global-foo"))

  resetGlobalEnv()
}
