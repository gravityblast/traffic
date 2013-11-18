package main

import (
  "github.com/pilu/traffic"
)

type ResponseData struct {
  Message string
}

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  responseData := &ResponseData{ "Hello World" }
  traffic.Render(w, "index", responseData)
}

func main() {
  router := traffic.New()
  router.Get("/", rootHandler)
  router.Run()
}
