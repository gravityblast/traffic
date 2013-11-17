package traffic

import (
  "net/http"
)

type ResponseWriter interface {
  http.ResponseWriter
  SetVar(string, interface{})
  GetVar(string) interface{}
  AddBeforeWriteHandler(handler func())
  StatusCode() int
  Written() bool
}

type AppResponseWriter struct {
  http.ResponseWriter
  written             bool
  statusCode          int
  env                 map[string]interface{}
  routerEnv           *map[string]interface{}
  beforeWriteHandlers []func()
}


func (w *AppResponseWriter) beforeWrite() {
  for _, handler := range w.beforeWriteHandlers {
    handler()
  }
}

func (w *AppResponseWriter) Write(data []byte) (n int, err error) {
  if !w.written {
    w.beforeWrite()
    w.written = true
  }

  return w.ResponseWriter.Write(data)
}

func (w *AppResponseWriter) WriteHeader(statusCode int) {
  if !w.written {
    w.beforeWrite()
  }

  w.statusCode = statusCode
  w.ResponseWriter.WriteHeader(statusCode)
  w.written = true
}

func (w *AppResponseWriter) StatusCode() int {
  return w.statusCode
}

func (w *AppResponseWriter) SetVar(key string, value interface{}) {
  w.env[key] = value
}

func (w *AppResponseWriter) Written() bool {
  return w.written
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

func (w *AppResponseWriter) AddBeforeWriteHandler(handler func()) {
  w.beforeWriteHandlers = append(w.beforeWriteHandlers, handler)
}

func newAppResponseWriter(w http.ResponseWriter, routerEnv *map[string]interface{}) *AppResponseWriter {
  arw := &AppResponseWriter{
    ResponseWriter:       w,
    statusCode:           http.StatusOK,
    env:                  make(map[string]interface{}),
    routerEnv:            routerEnv,
    beforeWriteHandlers:  make([]func(), 0),
  }

  return arw
}

