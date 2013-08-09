package traffic

import (
  "fmt"
  "strings"
)

func pathSegmentToRegexpSegment(segment string) string {
  var regexpSegment string
  if len(segment) < 1 {
    return segment
  } else if segment[0] == ':' {
    name := segment[1:len(segment)]
    regexpSegment = fmt.Sprintf("(?P<%s>[^/#?]+)", name)
  } else {
    regexpSegment = segment
  }

  return regexpSegment
}

func pathToRegexpString(path string) string {
  segments := strings.Split(path, "/")
  regexpSegments := make([]string, len(segments))

  for i, segment := range segments {
    regexpSegments[i] = pathSegmentToRegexpSegment(segment)
  }

  regexpString := strings.Join(regexpSegments, "/")
  return fmt.Sprintf("^%s$", regexpString)
}
