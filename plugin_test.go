package traefik_plugin_test

import (
	"context"
	"github.com/dragosnutu/traefik-plugin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlugin(t *testing.T) {
	cfg := traefik_plugin.CreateConfig()
	cfg.Seconds = 15
	cfg.Redirect = "some redirect"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_plugin.New(ctx, next, cfg, "plugin-demo")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, recorder, "X-Seconds", "15")
	assertHeader(t, recorder, "X-Redirect", "some redirect")
}

func assertHeader(t *testing.T, res *httptest.ResponseRecorder, key, expected string) {
	t.Helper()

	if res.Header().Get(key) != expected {
		t.Errorf("invalid header value: %s", res.Header().Get(key))
	}
}
