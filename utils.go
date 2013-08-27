package traffic

import (
  "fmt"
  "regexp"
)

func pathToRegexpString(path string) string {
  re := regexp.MustCompile(":[^/#?()]+")
  regexpString := re.ReplaceAllStringFunc(path, func(m string) string {
    return fmt.Sprintf("(?P<%s>[^/#?]+)", m[1:len(m)])
  })

  return fmt.Sprintf("^%s$", regexpString)
}
