package main

import (
  "net/http"
  "github.com/pilu/traffic"
  "fmt"
  "time"
)

func rootHandler(w traffic.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello World\n")
}

func privatePageHandler(w traffic.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello Private Page\n")
}

func pageHandler(w traffic.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  fmt.Fprintf(w, "Category ID: %s\n", params.Get("category_id"))
  fmt.Fprintf(w, "Page ID: %s\n", params.Get("id"))
}

func checkApiKey(w traffic.ResponseWriter, r *http.Request) bool {
  params := r.URL.Query()
  if params.Get("api_key") != "foo" {
    w.WriteHeader(http.StatusUnauthorized)
    return false
  }

  return true
}

func checkPrivatePageApiKey(w traffic.ResponseWriter, r *http.Request) bool {
  params := r.URL.Query()
  if params.Get("private_api_key") != "bar" {
    w.WriteHeader(http.StatusUnauthorized)
    return false
  }

  return true
}

func addAppNameHeader(w traffic.ResponseWriter, r *http.Request) bool {
  w.Header().Add("X-APP-NAME", "My App")

  return true
}

func addTimeHeader(w traffic.ResponseWriter, r *http.Request) bool {
  t := fmt.Sprintf("%s", time.Now())
  w.Header().Add("X-APP-TIME", t)

  return true
}

func main() {
  router := traffic.New()

  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)
  // "/private" has one more before filter that checks for a second api key (private_api_key)
  router.Get("/private", privatePageHandler).
          AddBeforeFilter(checkPrivatePageApiKey)

  // Executed before all handlers
  router.AddBeforeFilter(checkApiKey).
         AddBeforeFilter(addAppNameHeader).
         AddBeforeFilter(addTimeHeader)

  http.Handle("/", router)
  http.ListenAndServe(":7000", nil)
}
