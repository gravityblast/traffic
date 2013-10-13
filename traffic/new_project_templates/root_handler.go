package main

import (
  "net/http"
  "github.com/pilu/traffic"
)

type ResponseData struct {
  Message string
}

func RootHandler(w traffic.ResponseWriter, r *http.Request) {
  responseData := &ResponseData{ "Hello World" }
  traffic.Render(w, "index", responseData)
}

