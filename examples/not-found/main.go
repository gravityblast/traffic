package main

import (
  "fmt"
  "net/http"
  "github.com/pilu/traffic"
)

func rootHandler(w traffic.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello World\n")
}

func pageHandler(w traffic.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  fmt.Fprintf(w, "Category ID: %s\n", params.Get("category_id"))
  fmt.Fprintf(w, "Page ID: %s\n", params.Get("id"))
}

func customNotFoundHandler(w traffic.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Page not found: %s\n", r.URL.Path)
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
