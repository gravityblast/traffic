package main

import (
  "github.com/pilu/traffic"
)

func errorHandler(w traffic.ResponseWriter, r *traffic.Request, err interface{}) {
  w.WriteText("This is a custom error page <br />\n")
  w.WriteText("Recovered from `%v`", err)
}

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
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
