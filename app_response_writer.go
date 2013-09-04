package traffic

import (
  "net/http"
)

type AppResponseWriter struct {
  http.ResponseWriter
  statusCode int
}

func (w *AppResponseWriter) WriteHeader(statusCode int) {
  w.statusCode = statusCode
  w.ResponseWriter.WriteHeader(statusCode)
}

func (w *AppResponseWriter) StatusCode() int {
  if w.statusCode == 0 {
    return http.StatusOK
  }

  return w.statusCode
}

