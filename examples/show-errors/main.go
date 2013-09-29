package main

import (
  "log"
  "net/http"
  "github.com/pilu/traffic"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
  x := 0
  // this will raise a 'runtime error: integer divide by zero'
  x = 1 / x
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
