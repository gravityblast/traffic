package main

import (
  "strings"
  "github.com/pilu/traffic"
)

type ResponseData struct {
  Message string
}

func indexHandler(w traffic.ResponseWriter, r *traffic.Request) {
  responseData := &ResponseData{ "hello world" }
  w.Render("index", responseData)
}

func aboutHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.Render("about")
}

func main() {
  traffic.TemplateFunc("upcase",    strings.ToUpper)
  traffic.TemplateFunc("downcase",  strings.ToLower)

  router := traffic.New()
  router.Get("/", indexHandler)
  router.Get("/about/?", aboutHandler)
  router.Run()
}
