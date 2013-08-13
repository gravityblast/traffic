package main

import (
  "net/http"
  "github.com/pilu/traffic"
  "fmt"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello World\n")
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  fmt.Fprintf(w, "Category ID: %s\n", params.Get("category_id"))
  fmt.Fprintf(w, "Page ID: %s\n", params.Get("id"))
}

func customNotFoundHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusNotFound)
  fmt.Fprintf(w, "Page not found: %s\n", r.URL.Path)
}

func main() {
  router := traffic.New()

  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)

  // Custom not found handler
  router.NotFoundHandler = customNotFoundHandler

  http.Handle("/", router)
  http.ListenAndServe(":7000", nil)
}
