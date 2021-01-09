package traefik_kuzzle_auth

import (
	"context"
	"fmt"
	"net/http"
)

// KuzzleAuth a plugin to use Kuzzle as authentication provider for Basic Auth Traefik middleware.
type KuzzleAuth struct {
	next   http.Handler
	name   string
	config *Config
}

// New created a new KuzzleBasicAuth plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if err := config.Kuzzle.ping(); err != nil {
		return nil, fmt.Errorf("Unable to reach Kuzzle server at %s: %v", config.Kuzzle.URL, err)
	}

	return &KuzzleAuth{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

func (ka *KuzzleAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	user, pass, ok := req.BasicAuth()

	if !ok {
		// No valid 'Authentication: Basic xxxx' header found in request
		rw.Header().Set("WWW-Authenticate", `Basic realm="`+ka.config.CustomRealm+`"`)
		http.Error(rw, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	if err := ka.config.Kuzzle.login(user, pass); err != nil {
		// Failed to login with provided user/pass
		rw.Header().Set("WWW-Authenticate", `Basic realm="`+ka.config.CustomRealm+`"`)
		http.Error(rw, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	if len(ka.config.Kuzzle.AllowedUsers) > 0 {
		// Allowed Users have been specified
		if err := ka.config.Kuzzle.checkUser(); err != nil {
			// Logged user do not be part of the configured Allowed Users
			rw.Header().Set("WWW-Authenticate", `Basic realm="`+ka.config.CustomRealm+`"`)
			http.Error(rw, "Unauthorized.", http.StatusUnauthorized)
			return
		}
	}

	ka.next.ServeHTTP(rw, req)
}
