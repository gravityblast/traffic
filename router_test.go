package traffic

import (
  "testing"
  "net/http"
  "reflect"
  assert "github.com/pilu/miniassert"
)

func TestNew(t *testing.T) {
  router := New()
  assert.Type(t, "*traffic.Router", router)
  assert.Type(t, "map[traffic.HttpMethod][]*traffic.Route", router.routes)

  emptyMap := make(map[HttpMethod][]*Route)
  assert.Equal(t, emptyMap, router.routes)
}

func TestAdd(t *testing.T) {
  router := New()
  assert.Equal(t, 0, len(router.routes["GET"]))
  route := router.Add(HttpMethod("GET"), "/", httpHandlerExample)
  assert.Type(t, "*traffic.Route", route)
  assert.Equal(t, 1, len(router.routes["GET"]))
}

func TestGet(t *testing.T) {
  router := New()
  assert.Equal(t, 0, len(router.routes["GET"]))
  assert.Equal(t, 0, len(router.routes["HEAD"]))
  route := router.Get("/", httpHandlerExample)
  assert.Type(t, "*traffic.Route", route)
  assert.Equal(t, 1, len(router.routes["GET"]))
  assert.Equal(t, 1, len(router.routes["HEAD"]))
  assert.Equal(t, router.routes["GET"][0], router.routes["HEAD"][0])
}

func TestPost(t *testing.T) {
  router := New()
  assert.Equal(t, 0, len(router.routes["POST"]))
  route := router.Post("/", httpHandlerExample)
  assert.Type(t, "*traffic.Route", route)
  assert.Equal(t, 1, len(router.routes["POST"]))
}

func TestDelete(t *testing.T) {
  router := New()
  assert.Equal(t, 0, len(router.routes["DELETE"]))
  route := router.Delete("/", httpHandlerExample)
  assert.Type(t, "*traffic.Route", route)
  assert.Equal(t, 1, len(router.routes["DELETE"]))
}

func TestPut(t *testing.T) {
  router := New()
  assert.Equal(t, 0, len(router.routes["PUT"]))
  route := router.Put("/", httpHandlerExample)
  assert.Type(t, "*traffic.Route", route)
  assert.Equal(t, 1, len(router.routes["PUT"]))
}

func TestPatch(t *testing.T) {
  router := New()
  assert.Equal(t, 0, len(router.routes["PATCH"]))
  route := router.Patch("/", httpHandlerExample)
  assert.Type(t, "*traffic.Route", route)
  assert.Equal(t, 1, len(router.routes["PATCH"]))
}

func TestAddBeforeFilter(t *testing.T) {
  router := New()
  assert.Equal(t, 0, len(router.beforeFilters))

  filterA := BeforeFilterFunc(func(w http.ResponseWriter, r *http.Request) bool { return true })
  filterB := BeforeFilterFunc(func(w http.ResponseWriter, r *http.Request) bool { return true })

  router.AddBeforeFilter(filterA)
  assert.Equal(t, 1, len(router.beforeFilters))
  router.AddBeforeFilter(filterB)
  assert.Equal(t, 2, len(router.beforeFilters))

  assert.Equal(t, reflect.ValueOf(filterA), reflect.ValueOf(router.beforeFilters[0]))
  assert.Equal(t, reflect.ValueOf(filterB), reflect.ValueOf(router.beforeFilters[1]))
}

