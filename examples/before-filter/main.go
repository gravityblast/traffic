package main

import (
  "net/http"
  "github.com/pilu/traffic"
  "fmt"
  "time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello World\n")
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  fmt.Fprintf(w, "Category ID: %s\n", params.Get("category_id"))
  fmt.Fprintf(w, "Page ID: %s\n", params.Get("id"))
}

func checkApiKey(w http.ResponseWriter, r *http.Request) bool {
  params := r.URL.Query()
  if params.Get("api_key") != "foo" {
    w.WriteHeader(http.StatusUnauthorized)
    return false
  }

  return true
}

func addAppNameHeader(w http.ResponseWriter, r *http.Request) bool {
  w.Header().Add("X-APP-NAME", "My App")

  return true
}

func addTimeHeader(w http.ResponseWriter, r *http.Request) bool {
  t := fmt.Sprintf("%s", time.Now())
  w.Header().Add("X-APP-TIME", t)

  return true
}

func main() {
  router := traffic.New()

  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)

  // Executed before all handlers
  router.AddBeforeFilter(checkApiKey).
         AddBeforeFilter(addAppNameHeader).
         AddBeforeFilter(addTimeHeader)

  http.Handle("/", router)
  http.ListenAndServe(":7000", nil)
}
