package main

import (
  "github.com/pilu/traffic"
)

type ResponseData struct {
  Message string
}

func RootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  responseData := &ResponseData{ "Hello World" }
  w.Render("index", responseData)
}

