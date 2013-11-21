package main

import (
  "github.com/pilu/traffic"
)

type ResponseData struct {
  PagePath  string
}

func pageHandler(w traffic.ResponseWriter, r *traffic.Request) {
  pagePath  := r.Param("page_path")

  responseData := &ResponseData{
    PagePath: pagePath,
  }

  w.Render("index", responseData)
}

func main() {
  router := traffic.New()
  router.Get("/:page_path*?", pageHandler)
  router.Run()
}
