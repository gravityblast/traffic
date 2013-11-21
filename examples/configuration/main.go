package main

import (
  "github.com/pilu/traffic"
)

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.WriteText("%s<br />", w.GetVar("foo"))

  // run with TRAFFIC_ENV=production to get the "bar" value
  // from the production section of the config file
  w.WriteText("%s<br />", w.GetVar("bar"))
}

func main() {
  router := traffic.New()
  // Routes
  router.Get("/", rootHandler)
  router.Run()
}
