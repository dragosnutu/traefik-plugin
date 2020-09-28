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
	cfg.JwtField = "sub"
	cfg.JwtValues = []string{"1234567890", "abcd"}
	cfg.Redirect = "http://localhost/redirected"

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

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

	handler.ServeHTTP(recorder, req)

	assertHeader(t, recorder, "Location", "http://localhost/redirected")
	assertStatus(t, recorder, http.StatusTemporaryRedirect)
}

func assertHeader(t *testing.T, res *httptest.ResponseRecorder, key, expected string) {
	t.Helper()

	if res.Header().Get(key) != expected {
		t.Errorf("invalid header value: %s", res.Header().Get(key))
	}
}

func assertStatus(t *testing.T, res *httptest.ResponseRecorder, expected int) {
	t.Helper()

	actual := res.Result().StatusCode

	if actual != expected {
		t.Errorf("Invalid response status, expected: %d, actual: %d", expected, actual)
	}
}
