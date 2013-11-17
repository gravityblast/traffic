package traffic

import (
  "net/http"
  "net/url"
  "regexp"
)

type HttpHandleFunc func(ResponseWriter, *http.Request)

type Route struct {
  Path string
  PathRegexp *regexp.Regexp
  Handler HttpHandleFunc
  beforeFilters []HttpHandleFunc
}

func (route *Route) AddBeforeFilter(beforeFilter HttpHandleFunc) *Route {
  route.beforeFilters = append(route.beforeFilters, beforeFilter)

  return route
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
      value := matches[0][i]
      if len(name) > 0 && len(value) > 0 {
        values.Set(name, value)
      }
    }

    return values, true
  }

  return values, false
}

