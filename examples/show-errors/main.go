package main

import (
  "net/http"
  "github.com/pilu/traffic"
)

func rootHandler(w traffic.ResponseWriter, r *http.Request) {
  x := 0
  // this will raise a 'runtime error: integer divide by zero'
  x = 1 / x
}

func main() {
  router := traffic.New()
  router.Get("/", rootHandler)
  router.Run()
}
