package main

import (
  "net/http"
  "github.com/pilu/traffic"
)

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.WriteText("Hello World\n")
}

func pageHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.WriteText("Category ID: %s\n", r.Param("category_id"))
  w.WriteText("Page ID: %s\n", r.Param("id"))
}

func customNotFoundHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.WriteHeader(http.StatusNotFound)
  w.WriteText("Page not found: %s\n", r.URL.Path)
}

func main() {
  router := traffic.New()

  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)

  // Custom not found handler
  router.NotFoundHandler = customNotFoundHandler

  router.Run()
}
