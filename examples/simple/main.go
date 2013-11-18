package main

import (
  "fmt"
  "github.com/pilu/traffic"
)

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  traffic.Logger().Print("Hello")
  fmt.Fprint(w, "Hello World\n")
}

func pageHandler(w traffic.ResponseWriter, r *traffic.Request) {
  fmt.Fprintf(w, "Category ID: %s\n", r.Param("category_id"))
  fmt.Fprintf(w, "Page ID: %s\n", r.Param("id"))
}

func main() {
  router := traffic.New()
  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)

  router.Run()
}
