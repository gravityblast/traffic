package main

import (
  "net/http"
  "github.com/pilu/traffic"
  "fmt"
  "log"
  "os"
)

type AppLogger struct {
  Name string
  *log.Logger
}

func (appLogger *AppLogger) requestLogFunc(statusCode int, r *http.Request) {
  appLogger.Printf("%s: %d - %s\n", appLogger.Name, statusCode, r.URL)
}

var logger *AppLogger

func rootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello World\n")
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  fmt.Fprintf(w, "Category ID: %s\n", params.Get("category_id"))
  fmt.Fprintf(w, "Page ID: %s\n", params.Get("id"))
}

func customNotFoundHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusNotFound)
  fmt.Fprintf(w, "Page not found: %s\n", r.URL.Path)
}

func customBeforeFilter(w http.ResponseWriter, r *http.Request) bool {
  w.Header().Add("X-APP-NAME", "My App")

  return true
}

func init() {
  logger = &AppLogger{
    Name: "Test App",
    Logger: log.New(os.Stderr, "", log.LstdFlags),
  }
}

func main() {
  router := traffic.New()

  // Logger
  router.RequestLogFunc = logger.requestLogFunc

  // Routes
  router.Get("/", rootHandler)
  router.Get("/categories/:category_id/pages/:id", pageHandler)

  // Custom not found handler
  router.NotFoundHandler = customNotFoundHandler

  // Executed before all handlers
  router.AddBeforeFilter(customBeforeFilter)

  http.Handle("/", router)
  http.ListenAndServe(":7000", nil)
}
