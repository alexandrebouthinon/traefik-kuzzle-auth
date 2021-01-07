package kuzzlebasicauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Kuzzle struct {
		Host string `json:"host,omitempty"`
		Port uint   `json:"port,omitempty"`
		Ssl  bool   `json:"ssl,omitempty"`
		URL  string
	} `json:"kuzzle,omitempty"`
	BasicAuth struct {
		User     string `json:"user,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"basic-auth,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// KuzzleBasicAuth a plugin to use Kuzzle as authentication provider for Basic Auth Traefik middleware.
type KuzzleBasicAuth struct {
	next   http.Handler
	name   string
	config *Config
}

func (c *Config) check() error {
	if c.Kuzzle.Host == "" {
		return fmt.Errorf("You need to set proper value for 'host' field in 'kuzzle' configuration part")
	}

	if c.BasicAuth.User == "" || c.BasicAuth.Password == "" {
		return fmt.Errorf("You need to set proper values for 'user' and 'password' fields in 'basic-auth' configuration part")
	}

	return nil
}

func pingKuzzle(c *Config) (string, error) {
	proto := "http"

	if c.Kuzzle.Ssl == true {
		proto = "https"
	}

	url := fmt.Sprintf("%s://%s:%d", proto, c.Kuzzle.Host, c.Kuzzle.Port)
	pingRoute := fmt.Sprintf("%s/_publicApi", url) // This is a Kuzzle API route always exposed publicly
	_, err := http.Get(pingRoute)

	return url, err
}

// New created a new KuzzleBasicAuth plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if err := config.check(); err != nil {
		return nil, fmt.Errorf("Error during configuration check: %v", err)
	}

	url, err := pingKuzzle(config)
	if err != nil {
		return nil, fmt.Errorf("Unable to reach Kuzzle server at %s: %v", url, err)
	}
	config.Kuzzle.URL = url

	return &KuzzleBasicAuth{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

func (k *KuzzleBasicAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	user, pass, _ := req.BasicAuth()

	if ok, err := k.loginToKuzzleIsSuccessful(user, pass); !ok {
		rw.WriteHeader(http.StatusUnauthorized)
		if err != nil {
			rw.Write([]byte(err.Error()))
		}
		return
	}

	req.SetBasicAuth(k.config.BasicAuth.User, k.config.BasicAuth.Password)
	k.next.ServeHTTP(rw, req)
}

func (k *KuzzleBasicAuth) loginToKuzzleIsSuccessful(user string, password string) (bool, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"username": user,
		"password": password,
	})

	loginRoute := fmt.Sprintf("%s/_login/local", k.config.Kuzzle.URL)

	if _, err := http.Post(loginRoute, "application/json", bytes.NewBuffer(reqBody)); err != nil {
		return false, fmt.Errorf("Authentication request send to %s failed: %v", k.config.Kuzzle.URL, err)
	}

	return true, nil
}
