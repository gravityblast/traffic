package main

import (
  "github.com/pilu/traffic"
)

type ResponseData struct {
  Message string
}

func indexHandler(w traffic.ResponseWriter, r *traffic.Request) {
  responseData := &ResponseData{ "Hello World" }
  w.RenderTemplate("index", responseData)
}

func aboutHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.RenderTemplate("about")
}

func main() {
  router := traffic.New()
  router.Get("/", indexHandler)
  router.Get("/about/?", aboutHandler)
  router.Run()
}
