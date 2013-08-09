package traffic

import (
  "net/http"
  "net/url"
  "regexp"
)

type HttpHandleFunc func(http.ResponseWriter, *http.Request)

type Route struct {
  Path string
  PathRegexp *regexp.Regexp
  Handler HttpHandleFunc
}

func NewRoute(path string, handler HttpHandleFunc) *Route {
  route := &Route{}
  route.Path = path
  route.Handler = handler
  route.PathRegexp = regexp.MustCompile(pathToRegexpString(path))
  return route
}

func (route Route) Match(path string) (url.Values, bool) {
  values := make(url.Values)

  matches := route.PathRegexp.FindAllStringSubmatch(path, -1)
  if matches != nil {
    names := route.PathRegexp.SubexpNames()
    for i := 1; i < len(names); i++ {
      name := names[i]
      values.Set(name, matches[0][i])
    }

    return values, true
  }

  return values, false
}

