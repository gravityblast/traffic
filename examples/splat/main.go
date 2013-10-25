package main

import (
  "net/http"
  "github.com/pilu/traffic"
)

type ResponseData struct {
  PagePath  string
}

func pageHandler(w traffic.ResponseWriter, r *http.Request) {
  params    := r.URL.Query()
  pagePath  := params.Get("page_path")

  responseData := &ResponseData{
    PagePath: pagePath,
  }
  traffic.Render(w, "index", responseData)
}

func main() {
  router := traffic.New()
  router.Get("/:page_path*?", pageHandler)
  router.Run()
}
