package main

import (
  "github.com/pilu/traffic"
)

var router *traffic.Router

func init() {
  router = traffic.New()
  router.Get("/", RootHandler)
}

func main() {
  router.Run()
}
