package main

import (
  "fmt"
  "time"
  "net/http"
  "github.com/pilu/traffic"
)

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.WriteText("Hello World\n")
}

func privatePageHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.WriteText("Hello Private Page\n")
}

func pageHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.WriteText("Category ID: %s\n", r.Param("category_id"))
  w.WriteText("Page ID: %s\n", r.Param("id"))
}

func checkApiKey(w traffic.ResponseWriter, r *traffic.Request) {
  if r.Param("api_key") != "foo" {
    w.WriteHeader(http.StatusUnauthorized)
    w.WriteText("Not authorized\n")
  }
}

func checkPrivatePageApiKey(w traffic.ResponseWriter, r *traffic.Request) {
  if r.Param("private_api_key") != "bar" {
    w.WriteHeader(http.StatusUnauthorized)
    w.WriteText("Not authorized\n")
  }
}

func addAppNameHeader(w traffic.ResponseWriter, r *traffic.Request) {
  w.Header().Add("X-APP-NAME", "My App")
}

func addTimeHeader(w traffic.ResponseWriter, r *traffic.Request) {
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
