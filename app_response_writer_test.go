package traffic

import (
  "testing"
  "reflect"
  "net/http"
  "net/http/httptest"
  assert "github.com/pilu/miniassert"
)

func buildTestAppResponseWriter(globalEnv *map[string]interface{}) (*AppResponseWriter, *httptest.ResponseRecorder) {
  recorder := httptest.NewRecorder()
  arw := newAppResponseWriter(
    recorder,
    globalEnv,
  )
  return arw, recorder
}

func newTestAppResponseWriter(globalEnv *map[string]interface{}) *AppResponseWriter {
  arw, _ := buildTestAppResponseWriter(globalEnv)

  return arw
}

func TestAppResponseWriter(t *testing.T) {
  routerEnv := make(map[string]interface{})
  arw := newTestAppResponseWriter(&routerEnv)
  assert.Equal(t, http.StatusOK, arw.statusCode)
  assert.Equal(t, 0, len(arw.env))
  assert.Equal(t, &routerEnv, arw.routerEnv)
  assert.Equal(t, 0, len(arw.beforeWriteHandlers))
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

func TestAppResponseWriter_AddBeforeWriteHandler(t *testing.T) {
  routerEnv := make(map[string]interface{})
  arw := newTestAppResponseWriter(&routerEnv)


  handler := func() {}
  arw.AddBeforeWriteHandler(handler)

  assert.Equal(t, 1, len(arw.beforeWriteHandlers))
  assert.Equal(t, reflect.ValueOf(handler), reflect.ValueOf(arw.beforeWriteHandlers[0]))
}

type beforeWriteHandlerTest struct {
  calls int
}

func (b *beforeWriteHandlerTest) handler() {
  b.calls++
}

func TestAppResponseWriter_Write(t *testing.T) {
  routerEnv := make(map[string]interface{})
  arw, recorder := buildTestAppResponseWriter(&routerEnv)

  handler_1 := &beforeWriteHandlerTest{}
  handler_2 := &beforeWriteHandlerTest{}
  arw.AddBeforeWriteHandler(handler_1.handler)
  arw.AddBeforeWriteHandler(handler_2.handler)

  assert.False(t, arw.Written())
  assert.Equal(t, 0, handler_1.calls)
  assert.Equal(t, 0, handler_2.calls)

  arw.Write([]byte("foo"))

  assert.True(t, arw.Written())
  assert.Equal(t, 1, handler_1.calls)
  assert.Equal(t, 1, handler_2.calls)
  assert.Equal(t, []byte("foo"), recorder.Body.Bytes())

  arw.Write([]byte("bar"))

  // handlers are not called the second time we call Write
  assert.Equal(t, 1, handler_1.calls)
  assert.Equal(t, 1, handler_2.calls)
  assert.Equal(t, []byte("foobar"), recorder.Body.Bytes())
}
