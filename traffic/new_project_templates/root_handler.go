package main

import (
  "github.com/pilu/traffic"
)

type ResponseData struct {
  Message string
}

func RootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  responseData := &ResponseData{ "Hello World" }
  traffic.Render(w, "index", responseData)
}

