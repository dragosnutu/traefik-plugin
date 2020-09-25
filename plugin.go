package traefik_plugin

import (
	"context"
	"fmt"
	"net/http"
)

type Config struct {
	Seconds  int32  `json:"seconds,omitempty"`
	Redirect string `json:"redirect"`
}

func CreateConfig() *Config {
	return &Config{Seconds: int32(60)}
}

type Plugin struct {
	next     http.Handler
	seconds  int32
	redirect string
	name     string
}

//goland:noinspection GoUnusedParameter
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Seconds < 0 {
		return nil, fmt.Errorf("Positive seconds are required")
	}

	if len(config.Redirect) == 0 {
		return nil, fmt.Errorf("Redirect config is required")
	}

	return &Plugin{
		seconds:  config.Seconds,
		redirect: config.Redirect,
		next:     next,
		name:     name,
	}, nil
}

func (a *Plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Add("X-Seconds", fmt.Sprint(a.seconds))
	rw.Header().Add("X-Redirect", a.redirect)

	a.next.ServeHTTP(rw, req)
}
