package main

import (
  "log"
  "net/http"
  "github.com/pilu/traffic"
)

type ResponseData struct {
  Message string
}

func rootHandler(w traffic.ResponseWriter, r *http.Request) {
  responseData := &ResponseData{ "Hello World" }
  traffic.Render(w, "index", responseData)
}

func main() {
  router := traffic.New()

  // Routes
  router.Get("/", rootHandler)

  http.Handle("/", router)
  err := http.ListenAndServe(":7000", nil)
  if err != nil {
    log.Fatal(err)
  }
}
