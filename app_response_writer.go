package traffic

import (
  "net/http"
)

type AppResponseWriter struct {
  http.ResponseWriter
  statusCode int
  env map[string]interface{}
  routerEnv *map[string]interface{}
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
  // local env
  value := w.env[key]
  if value != nil {
    return value
  }

  // router env
  value = (*w.routerEnv)[key]
  if value != nil {
    return value
  }

  // global env
  return GetVar(key)
}

func newAppResponseWriter(w http.ResponseWriter, routerEnv *map[string]interface{}) *AppResponseWriter {
  arw := &AppResponseWriter{
    w,
    http.StatusOK,
    make(map[string]interface{}),
    routerEnv,
  }

  return arw
}

