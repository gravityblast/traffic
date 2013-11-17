package main

import (
  "fmt"
  "time"
  "net/http"
  "github.com/pilu/traffic"
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

func checkApiKey(w traffic.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  if params.Get("api_key") != "foo" {
    w.WriteHeader(http.StatusUnauthorized)
  }
}

func checkPrivatePageApiKey(w traffic.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  if params.Get("private_api_key") != "bar" {
    w.WriteHeader(http.StatusUnauthorized)
  }
}

func addAppNameHeader(w traffic.ResponseWriter, r *http.Request) {
  w.Header().Add("X-APP-NAME", "My App")
}

func addTimeHeader(w traffic.ResponseWriter, r *http.Request) {
  t := fmt.Sprintf("%s", time.Now())
  w.Header().Add("X-APP-TIME", t)
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
  router.AddBeforeFilter(checkApiKey, addAppNameHeader, addTimeHeader)

  router.Run()
}
