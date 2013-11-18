package traffic

import (
  "net/http"
)

type ResponseWriter interface {
  http.ResponseWriter
  SetVar(string, interface{})
  GetVar(string) interface{}
  StatusCode() int
  Written() bool
}

type responseWriter struct {
  http.ResponseWriter
  written             bool
  statusCode          int
  env                 map[string]interface{}
  routerEnv           *map[string]interface{}
  beforeWriteHandlers []func()
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
  if !w.written {
    w.written = true
  }

  return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteHeader(statusCode int) {
  w.statusCode = statusCode
  w.ResponseWriter.WriteHeader(statusCode)
  w.written = true
}

func (w *responseWriter) StatusCode() int {
  return w.statusCode
}

func (w *responseWriter) SetVar(key string, value interface{}) {
  w.env[key] = value
}

func (w *responseWriter) Written() bool {
  return w.written
}

func (w *responseWriter) GetVar(key string) interface{} {
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

func newResponseWriter(w http.ResponseWriter, routerEnv *map[string]interface{}) *responseWriter {
  rw := &responseWriter{
    ResponseWriter:       w,
    statusCode:           http.StatusOK,
    env:                  make(map[string]interface{}),
    routerEnv:            routerEnv,
    beforeWriteHandlers:  make([]func(), 0),
  }

  return rw
}

