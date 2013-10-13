package main

import (
  "log"
  "net/http"
  "net/http/httptest"
)

func newTestRequest(method, path string) *httptest.ResponseRecorder {
  request, err := http.NewRequest(method, path, nil)
  if err != nil {
    log.Fatal(err)
  }

  recorder := httptest.NewRecorder()
  router.ServeHTTP(recorder, request)

  return recorder
}
