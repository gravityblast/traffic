package main

import (
  "net/http"
  "github.com/pilu/traffic"
  "fmt"
)

func rootHandler(r http.ResponseWriter, req *http.Request) {
  fmt.Fprint(r, "Hello World")
}

func pageHandler(r http.ResponseWriter, req *http.Request) {
  params := req.URL.Query()
  fmt.Fprintf(r, "Category ID: %s\n", params.Get("category_id"))
  fmt.Fprintf(r, "Page ID: %s\n", params.Get("id"))
}

func main() {
  router := traffic.New()
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)
  http.Handle("/", router)
  http.ListenAndServe(":7000", nil)
}
