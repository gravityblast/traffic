package main

import (
  "fmt"
  "net/http"
  "github.com/pilu/traffic"
)

func rootHandler(w traffic.ResponseWriter, r *http.Request) {
  traffic.Logger().Print("Hello")
  fmt.Fprint(w, "Hello World\n")
}

func pageHandler(w traffic.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  fmt.Fprintf(w, "Category ID: %s\n", params.Get("category_id"))
  fmt.Fprintf(w, "Page ID: %s\n", params.Get("id"))
}

func main() {
  router := traffic.New()
  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)

  router.Run()
}
