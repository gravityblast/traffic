package main

import (
  "fmt"
  "net/http"
  "github.com/pilu/traffic"
)

func errorHandler(w traffic.ResponseWriter, r *http.Request, err interface{}) {
  fmt.Fprint(w, "This is a custom error page <br />\n")
  fmt.Fprintf(w, "Recovered from `%v`", err)
}

func rootHandler(w traffic.ResponseWriter, r *http.Request) {
  x := 0
  // this will raise a 'runtime error: integer divide by zero'
  x = 1 / x
}

func main() {
  // Setting env to `production`.
  // Otherwise the ShoeErrors Middleware used in development
  // will recover from the panics.

  traffic.SetVar("env", "production")

  router := traffic.New()
  router.ErrorHandler = errorHandler
  router.Get("/", rootHandler)
  router.Run()
}
