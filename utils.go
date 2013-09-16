package traffic

import (
  "fmt"
  "regexp"
  "os"
  "strings"
)

var env map[string]interface{}

func init() {
  env = make(map[string]interface{})
}

func SetVar(key string, value interface{}) {
  env[key] = value
}

func GetVar(key string) interface{} {
  if value := env[key]; value != nil {
    return value
  }

  envKey := fmt.Sprintf("TRAFFIC_%s", strings.ToUpper(key))
  if value := os.Getenv(envKey); value != "" {
    return value
  }

  return nil
}

func pathToRegexpString(path string) string {
  re := regexp.MustCompile(":[^/#?()]+")
  regexpString := re.ReplaceAllStringFunc(path, func(m string) string {
    return fmt.Sprintf("(?P<%s>[^/#?]+)", m[1:len(m)])
  })

  return fmt.Sprintf("^%s$", regexpString)
}
