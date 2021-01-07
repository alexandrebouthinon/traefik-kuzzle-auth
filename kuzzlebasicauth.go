package kuzzle_basic_auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	Kuzzle struct {
		Host string `json:"host,omitempty"`
		Port uint   `json:"port,omitempty"`
		Ssl  bool   `json:"ssl,omitempty"`
		Url  string
	} `json:"kuzzle,omitempty"`
	BasicAuth struct {
		User     string `json:"user,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"basic-auth,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type KuzzleBasicAuth struct {
	next   http.Handler
	name   string
	config *Config
}

func (c *Config) Check() error {
	if c.Kuzzle.Host == "" {
		return fmt.Errorf("You need to set proper value for 'host' field in 'kuzzle' configuration part")
	}

	if c.BasicAuth.User == "" || c.BasicAuth.Password == "" {
		return fmt.Errorf("You need to set proper values for 'user' and 'password' fields in 'basic-auth' configuration part")
	}

	return nil
}

func PingKuzzle(c *Config) (string, error) {
	proto := "http"

	if c.Kuzzle.Ssl == true {
		proto = "https"
	}

	url := fmt.Sprintf("%s://%s:%d", proto, c.Kuzzle.Host, c.Kuzzle.Port)
	pingRoute := fmt.Sprintf("%s/_publicApi", url) // This is a Kuzzle API route always exposed publicly
	_, err := http.Get(pingRoute)

	return url, err
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if err := config.Check(); err != nil {
		return nil, fmt.Errorf("Error during configuration check: %v", err)
	}

	url, err := PingKuzzle(config)
	if err != nil {
		return nil, fmt.Errorf("Unable to reach Kuzzle server at %s: %v", url, err)
	}
	config.Kuzzle.Url = url

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

	loginRoute := fmt.Sprintf("%s/_login/local", k.config.Kuzzle.Url)

	if _, err := http.Post(loginRoute, "application/json", bytes.NewBuffer(reqBody)); err != nil {
		return false, fmt.Errorf("Authentication request send to %s failed: %v", k.config.Kuzzle.Url, err)
	}

	return true, nil
}
