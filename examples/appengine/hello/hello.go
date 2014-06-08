package hello

// UNCOMMENT LINE 6 AND 12-14 TO RUN THIS EXAMPLE. THEY ARE COMMENTED TO RUN ON travis-ci.org WITHOUT ERRORS

import (
  /* "appengine" */
  "net/http"
  "github.com/pilu/traffic"
)

func init() {
  /* if !appengine.IsDevAppServer() { */
  /*   traffic.SetVar("env", "production") */
  /* } */

  t := traffic.New()
  t.Get("/", rootHandler)

  http.Handle("/", t)
}

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.Render("index", struct{
    Message string
  }{
    "Hello Google App Engine",
  })
}
