package main

import (
  "github.com/pilu/traffic"
)

type ResponseData struct {
  Message string
}

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  responseData := &ResponseData{ "Hello World" }
  w.Render("index", responseData)
}

func main() {
  router := traffic.New()
  router.Get("/", rootHandler)
  router.Run()
}
