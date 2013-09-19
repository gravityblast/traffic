package hello

import (
  "appengine"
  "net/http"
  "github.com/pilu/traffic"
)

func init() {
  if !appengine.IsDevAppServer() {
    traffic.SetVar("env", "production")
  }

  t := traffic.New()
  t.Get("/", rootHandler)

  http.Handle("/", t)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
  traffic.Render(w, "index", struct{
    Message string
  }{
    "Hello Google App Engine",
  })
}
