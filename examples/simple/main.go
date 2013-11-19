package main

import (
  "fmt"
  "github.com/pilu/traffic"
)

func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
  traffic.Logger().Print("Hello")
  fmt.Fprint(w, "Hello World\n")
}

func jsonTestHandler(w traffic.ResponseWriter, r *traffic.Request) {
  data := map[string]string{
    "foo": "bar",
  }
  w.RenderJSON(data)
}

func xmlTestHandler(w traffic.ResponseWriter, r *traffic.Request) {
  type Person struct {
    FirstName string   `xml:"name>first"`
    LastName  string   `xml:"name>last"`
  }

  w.RenderXML(&Person{
    FirstName:  "foo",
    LastName:   "bar",
  })
}

func pageHandler(w traffic.ResponseWriter, r *traffic.Request) {
  w.RenderText("Category ID: %s\n", r.Param("category_id"))
  w.RenderText("Page ID: %s\n", r.Param("id"))
}

func main() {
  router := traffic.New()
  // Routes
  router.Get("/",     rootHandler)
  router.Get("/json", jsonTestHandler)
  router.Get("/xml",  xmlTestHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)

  router.Run()
}
