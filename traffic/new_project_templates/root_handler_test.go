package main

import (
  "testing"
)

func TestRootHandler(t *testing.T) {
  recorder := newTestRequest("GET", "/")
  expectedStatusCode := 200
  if recorder.Code != expectedStatusCode {
    t.Errorf("Expected response status code `%d`, got `%d`", expectedStatusCode, recorder.Code)
  }
}
