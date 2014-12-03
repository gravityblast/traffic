package traffic

import (
	"net/url"
	"regexp"
)

type Route struct {
	Path          string
	PathRegexp    *regexp.Regexp
	IsStatic      bool
	Handlers      []HttpHandleFunc
	beforeFilters []HttpHandleFunc
}

func (route *Route) AddBeforeFilter(beforeFilters ...HttpHandleFunc) *Route {
	route.beforeFilters = append(route.beforeFilters, beforeFilters...)

	return route
}

func NewRoute(path string, handlers ...HttpHandleFunc) *Route {
	regexp, isStatic := pathToRegexp(path)
	route := &Route{
		Path:       path,
		PathRegexp: regexp,
		IsStatic:   isStatic,
		Handlers:   handlers,
	}

	return route
}

func (route Route) Match(path string) (url.Values, bool) {
	values := make(url.Values)

	if route.IsStatic {
		return values, route.Path == path
	}

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
