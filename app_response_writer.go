package traffic

import (
  "net/http"
)

type AppResponseWriter struct {
  http.ResponseWriter
  statusCode int
  env map[string]interface{}
  globalEnv *map[string]interface{}
}

func (w *AppResponseWriter) WriteHeader(statusCode int) {
  w.statusCode = statusCode
  w.ResponseWriter.WriteHeader(statusCode)
}

func (w *AppResponseWriter) StatusCode() int {
  return w.statusCode
}

func (w *AppResponseWriter) SetVar(key string, value interface{}) {
  w.env[key] = value
}

func (w *AppResponseWriter) GetVar(key string) interface{} {
  value := w.env[key]
  if value != nil {
    return value
  }

  return (*w.globalEnv)[key]
}

func newAppResponseWriter(w http.ResponseWriter, globalEnv *map[string]interface{}) *AppResponseWriter {
  arw := &AppResponseWriter{
    w,
    http.StatusOK,
    make(map[string]interface{}),
    globalEnv,
  }

  return arw
}

