package traffic

import (
	assert "github.com/pilu/miniassert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func buildTestResponseWriter(globalEnv *map[string]interface{}) (*responseWriter, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	rw := newResponseWriter(
		recorder,
		globalEnv,
	)
	return rw, recorder
}

func newTestResponseWriter(globalEnv *map[string]interface{}) *responseWriter {
	rw, _ := buildTestResponseWriter(globalEnv)

	return rw
}

func TestResponseWriter(t *testing.T) {
	routerEnv := make(map[string]interface{})
	rw := newTestResponseWriter(&routerEnv)
	assert.Equal(t, http.StatusOK, rw.statusCode)
	assert.Equal(t, 0, len(rw.env))
	assert.Equal(t, &routerEnv, rw.routerEnv)
	assert.Equal(t, 0, len(rw.beforeWriteHandlers))
}

func TestResponseWriter_SetVar(t *testing.T) {
	globalEnv := make(map[string]interface{})
	rw := newTestResponseWriter(&globalEnv)
	rw.SetVar("foo", "bar")
	assert.Equal(t, "bar", rw.env["foo"])
}

func TestResponseWriter_GetVar(t *testing.T) {
	resetGlobalEnv()

	routerEnv := map[string]interface{}{}
	rw := newTestResponseWriter(&routerEnv)

	env["global-foo"] = "global-bar"
	assert.Equal(t, "global-bar", rw.GetVar("global-foo"))

	routerEnv["global-foo"] = "router-bar"
	assert.Equal(t, "router-bar", rw.GetVar("global-foo"))

	rw.env["global-foo"] = "local-bar"
	assert.Equal(t, "local-bar", rw.GetVar("global-foo"))

	resetGlobalEnv()
}

func TestResponseWriter_Write(t *testing.T) {
	routerEnv := make(map[string]interface{})
	rw, recorder := buildTestResponseWriter(&routerEnv)

	assert.False(t, rw.Written())

	rw.Write([]byte("foo"))

	assert.True(t, rw.Written())
	assert.Equal(t, []byte("foo"), recorder.Body.Bytes())
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	routerEnv := make(map[string]interface{})
	rw, recorder := buildTestResponseWriter(&routerEnv)

	assert.False(t, rw.Written())

	rw.WriteHeader(http.StatusUnauthorized)

	assert.True(t, rw.Written())
	assert.Equal(t, http.StatusUnauthorized, rw.StatusCode())
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
