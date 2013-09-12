package main

import (
  "fmt"
  "github.com/pilu/traffic"
  "net/http"
  "log"
)

type PingMiddleware struct {}

// If the request path is "/ping", it writes PONG in the response and returns without calling the next middleware
// Otherwise it sets the variable "PING" with PONG as value and calls the next  middleware.
// The next middleware can
func (c *PingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next traffic.NextMiddlewareFunc) (http.ResponseWriter, *http.Request) {
  if r.URL.Path == "/ping" {
    fmt.Fprint(w, "pong\n")

    return w, r
  }

  if nextMiddleware := next(); nextMiddleware != nil {
    arw := w.(*traffic.AppResponseWriter)
    arw.SetVar("ping", "pong")
    w, r = nextMiddleware.ServeHTTP(w, r, next)
  }

  return w, r
}

func root(w http.ResponseWriter, r *http.Request) {
  arw := w.(*traffic.AppResponseWriter)

  logger := arw.GetVar("logger").(*log.Logger)
  logger.Printf("Hello")

  fmt.Fprintf(w, "Global var foo: %v\n", arw.GetVar("foo"))
  fmt.Fprintf(w, "Middleware var PING: %v\n", arw.GetVar("ping"))
}

func main() {
  t := traffic.New()
  // Add PingMiddleware
  t.AddMiddleware(&PingMiddleware{})
  // Set global var "foo"
  t.SetVar("foo", "bar")
  // Add root handler
  t.Get("/", root)

  http.Handle("/", t)
  http.ListenAndServe(":7000", nil)
}
