package traffic

import (
  "testing"
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
  router.Add(HttpMethod("GET"), &Route{})
  assert.Equal(t, 1, len(router.routes["GET"]))
}

func TestGet(t *testing.T) {
  router := New()
  assert.Equal(t, 0, len(router.routes["GET"]))
  assert.Equal(t, 0, len(router.routes["HEAD"]))
  router.Get("/", httpHandlerExample)
  assert.Equal(t, 1, len(router.routes["GET"]))
  assert.Equal(t, 1, len(router.routes["HEAD"]))
}
