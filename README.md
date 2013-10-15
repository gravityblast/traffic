# Traffic

[![Build Status](https://travis-ci.org/pilu/traffic.png?branch=master)](https://travis-ci.org/pilu/traffic)

Package traffic - a Sinatra inspired regexp/pattern mux for [Go](http://golang.org/ "The Go programming language").
[Mailing List](https://groups.google.com/forum/#!forum/traffic-go)

## Installation

    go get github.com/pilu/traffic

## Features

  * [Regexp routing](https://github.com/pilu/traffic/blob/master/examples/simple/main.go)
  * [Before Filters](https://github.com/pilu/traffic/blob/master/examples/before-filter/main.go)
  * [Custom not found handler](https://github.com/pilu/traffic/blob/master/examples/not-found/main.go)
  * [Middlewares](https://github.com/pilu/traffic/blob/master/examples/middleware/main.go)
  * [Templates/Views](https://github.com/pilu/traffic/tree/master/examples/templates)
  * [Easy Configuration](https://github.com/pilu/traffic/tree/master/examples/configuration)

## Development Features

  * [Shows errors and stacktrace in browser](https://github.com/pilu/traffic/tree/master/examples/show-errors)
  * [Serves static files](https://github.com/pilu/traffic/tree/master/examples/static-files)
  * Project Generator

## Installation

Dowload the `Traffic` code:

```bash
go get github.com/pilu/traffic
```

Build the command line tool:

```bash
go get github.com/pilu/traffic/traffic
```

Create a new project:
```bash
traffic new hello
```

Run your project:
```bash
cd hello
go build && ./hello
```

A command line `runner` is in development. It will build and run the project every time a file is created/modified.

## Example:
The following code is a simple example, the documentation in still in development.
For more examples check the `examples` folder.

```go
package main

import (
  "net/http"
  "github.com/pilu/traffic"
  "fmt"
)

func rootHandler(w traffic.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello World\n")
}

func pageHandler(w traffic.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  fmt.Fprintf(w, "Category ID: %s\n", params.Get("category_id"))
  fmt.Fprintf(w, "Page ID: %s\n", params.Get("id"))
}

func main() {
  router := traffic.New()

  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)

  http.Handle("/", router)
  http.ListenAndServe(":7000", nil)
}
```

## Before Filters

You can also add "before filters" to all your routes or just to some of them:

```go
router := traffic.New()

// Executed before all handlers
router.AddBeforeFilter(checkApiKey).
       AddBeforeFilter(addAppNameHeader).
       AddBeforeFilter(addTimeHeader)

// Routes
router.Get("/", rootHandler)
router.Get("/categories/:category_id/pages/:id", pageHandler)

// "/private" has one more before filter that checks for a second api key (private_api_key)
router.Get("/private", privatePageHandler).
        AddBeforeFilter(checkPrivatePageApiKey)
```

Complete example:

```go
func rootHandler(w traffic.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello World\n")
}

func privatePageHandler(w traffic.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello Private Page\n")
}

func pageHandler(w traffic.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  fmt.Fprintf(w, "Category ID: %s\n", params.Get("category_id"))
  fmt.Fprintf(w, "Page ID: %s\n", params.Get("id"))
}

func checkApiKey(w traffic.ResponseWriter, r *http.Request) bool {
  params := r.URL.Query()
  if params.Get("api_key") != "foo" {
    w.WriteHeader(http.StatusUnauthorized)
    return false
  }

  return true
}

func checkPrivatePageApiKey(w traffic.ResponseWriter, r *http.Request) bool {
  params := r.URL.Query()
  if params.Get("private_api_key") != "bar" {
    w.WriteHeader(http.StatusUnauthorized)
    return false
  }

  return true
}

func addAppNameHeader(w traffic.ResponseWriter, r *http.Request) bool {
  w.Header().Add("X-APP-NAME", "My App")

  return true
}

func addTimeHeader(w traffic.ResponseWriter, r *http.Request) bool {
  t := fmt.Sprintf("%s", time.Now())
  w.Header().Add("X-APP-TIME", t)

  return true
}

func main() {
  router := traffic.New()

  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)
  // "/private" has one more before filter that checks for a second api key (private_api_key)
  router.Get("/private", privatePageHandler).
          AddBeforeFilter(checkPrivatePageApiKey)

  // Executed before all handlers
  router.AddBeforeFilter(checkApiKey).
         AddBeforeFilter(addAppNameHeader).
         AddBeforeFilter(addTimeHeader)

  http.Handle("/", router)
  http.ListenAndServe(":7000", nil)
}
```

## Author

* [Andrea Franz](http://gravityblast.com)
