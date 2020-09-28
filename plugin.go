package traefik_plugin

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Config struct {
	JwtField  string   `json:"jwtFields"`
	JwtValues []string `json:"jwtValues"`
	Redirect  string   `json:"redirect"`
}

func CreateConfig() *Config {
	return &Config{}
}

type Plugin struct {
	name      string
	next      http.Handler
	jwtField  string
	jwtValues []string
	redirect  string
}

type Token struct {
	header  string
	payload string
	sign    string
}

//goland:noinspection GoUnusedParameter
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	fmt.Printf("Configs, jwtField=%v , jwtValues=%v, redirect=%v\n", config.JwtField, config.JwtValues, config.Redirect)

	if len(config.JwtField) == 0 {
		return nil, fmt.Errorf("jwtField needs to be set, current value=%v", config.JwtField)
	}

	if len(config.JwtValues) == 0 {
		return nil, fmt.Errorf("jwtValues needs to be set, current values=%v", config.JwtValues)
	}

	if len(config.Redirect) == 0 {
		return nil, fmt.Errorf("redirect needs to be set, current value=%v", config.Redirect)
	}

	return &Plugin{
		name:      name,
		next:      next,
		jwtField:  config.JwtField,
		jwtValues: config.JwtValues,
		redirect:  config.Redirect,
	}, nil
}

func (a *Plugin) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	fmt.Printf("Using auth value=%v\n", authHeader)

	if len(authHeader) != 0 {

		jwtValue := strings.TrimPrefix(authHeader, "Bearer ")
		jwtValue = strings.TrimSpace(jwtValue)

		parts := strings.Split(jwtValue, ".")

		var token Token

		token.header = parts[0]
		token.payload = parts[1]
		token.sign = parts[2]

		jwtDecodedValue, err := base64.RawURLEncoding.DecodeString(token.payload)
		if err != nil {
			http.Error(res, "We could not decode jwt value", http.StatusBadRequest)
		}

		fmt.Printf("Jwt decode value=%v\n", string(jwtDecodedValue))

		var rawJwt map[string]interface{}

		err = json.Unmarshal(jwtDecodedValue, &rawJwt)
		fmt.Printf("Unmarshalled raw %+v\n", rawJwt)
		if err != nil {
			http.Error(res, "We could not extract values", http.StatusBadRequest)
		}

		jwtFieldValue := rawJwt[a.jwtField].(string)

		fmt.Printf("Checking jwtFieldValue=%v to jwtValues=%v\n", jwtFieldValue, a.jwtValues)
		if contains(a.jwtValues, jwtFieldValue) {
			fmt.Printf("We have %v on jwtField=%v, and jwt_values=%v contains it so -> redirecting to %v\n", jwtFieldValue, a.jwtField, a.jwtValues, a.redirect)
			res.Header().Add("Location", a.redirect)
			res.WriteHeader(http.StatusTemporaryRedirect)

			return
		}
	}

	a.next.ServeHTTP(res, req)
}

func contains(array []string, val string) bool {
	for _, el := range array {
		equal := el == val
		fmt.Printf("Checking if val=%+v %T to el=%+v %T are same=%+v\n", val, val, el, el, equal)
		if equal {
			return true
		}
	}

	return false
}
