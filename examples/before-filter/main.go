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

func customBeforeFilter(w http.ResponseWriter, r *http.Request) bool {
  params := r.URL.Query()
  if params.Get("api_key") != "foo" {
    w.WriteHeader(http.StatusUnauthorized)
    return false
  }
  return true
}

func main() {
  router := traffic.New()

  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)

  // Executed before all handlers
  router.BeforeFilter = customBeforeFilter

  http.Handle("/", router)
  http.ListenAndServe(":7000", nil)
}
