package traefik_kuzzle_auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func TestNew(t *testing.T) {
	type args struct {
		ctx    context.Context
		next   http.Handler
		config *Config
		name   string
	}
	tests := []struct {
		name    string
		args    args
		want    http.Handler
		wantErr bool
		mock    Mock
	}{
		{
			name: "Valid configuration",
			args: args{
				config: &Config{
					Kuzzle: Kuzzle{
						URL:    "http://kuzzle:7512",
						Routes: Routes{},
					},
				},
			},
			want: &KuzzleAuth{
				config: &Config{
					Kuzzle: Kuzzle{
						URL:    "http://kuzzle:7512",
						Routes: Routes{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock.enabled {
				defer gock.Off()
				gock.
					New(tt.mock.url).
					Get(tt.mock.route).
					Reply(tt.mock.statusCode)
			}
			got, err := New(tt.args.ctx, tt.args.next, tt.args.config, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKuzzleAuth_ServeHTTP(t *testing.T) {
	type fields struct {
		name   string
		config *Config
	}
	tests := []struct {
		name               string
		fields             fields
		expNextCall        bool
		expStatusCode      int
		loginMock          Mock
		getCurrentUserMock Mock
		basicAuthMock      struct {
			enabled  bool
			username string
			password string
		}
	}{
		{
			name: "No Basic Auth info provided",
			fields: fields{
				name: "kuzzle-auth",
				config: &Config{
					Kuzzle: Kuzzle{
						URL: "http://kuzzle:7512",
						Routes: Routes{
							Login:          "/_login/local",
							GetCurrentUser: "/_me",
						},
					},
					CustomRealm: "Use a valid Kuzzle user to authenticate",
				},
			},
			expNextCall:   false,
			expStatusCode: 401,
		},
		{
			name: "Fail to log in to Kuzzle",
			fields: fields{
				name: "kuzzle-auth",
				config: &Config{
					Kuzzle: Kuzzle{
						URL: "http://kuzzle:7512",
						Routes: Routes{
							Login:          "/_login/local",
							GetCurrentUser: "/_me",
						},
					},
					CustomRealm: "Use a valid Kuzzle user to authenticate",
				},
			},
			expNextCall:   false,
			expStatusCode: 401,
			basicAuthMock: struct {
				enabled  bool
				username string
				password string
			}{
				enabled:  true,
				username: "test",
				password: "test",
			},
		},
		{
			name: "Log in succeed but get current user fail",
			fields: fields{
				name: "kuzzle-auth",
				config: &Config{
					Kuzzle: Kuzzle{
						URL: "http://kuzzle:7512",
						Routes: Routes{
							Login:          "/_login/local",
							GetCurrentUser: "/_me",
						},
						AllowedUsers: []string{"test"},
					},
					CustomRealm: "Use a valid Kuzzle user to authenticate",
				},
			},
			expNextCall:   false,
			expStatusCode: 401,
			basicAuthMock: struct {
				enabled  bool
				username string
				password string
			}{
				enabled:  true,
				username: "test",
				password: "test",
			},
			loginMock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_login/local",
				response:   json.RawMessage(`{"result":{"jwt": "myToken"}}`),
			},
		},
		{
			name: "Log in succeed but given user is not allowed",
			fields: fields{
				name: "kuzzle-auth",
				config: &Config{
					Kuzzle: Kuzzle{
						URL: "http://kuzzle:7512",
						Routes: Routes{
							Login:          "/_login/local",
							GetCurrentUser: "/_me",
						},
						AllowedUsers: []string{"toto"},
					},
					CustomRealm: "Use a valid Kuzzle user to authenticate",
				},
			},
			expNextCall:   false,
			expStatusCode: 401,
			basicAuthMock: struct {
				enabled  bool
				username string
				password string
			}{
				enabled:  true,
				username: "test",
				password: "test",
			},
			loginMock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_login/local",
				response:   json.RawMessage(`{"result":{"jwt": "myToken"}}`),
			},
			getCurrentUserMock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_login/local",
				response:   json.RawMessage(`{"result":{"_id": "test"}}`),
			},
		},
		{
			name: "Success with allowedUsers configured",
			fields: fields{
				name: "kuzzle-auth",
				config: &Config{
					Kuzzle: Kuzzle{
						URL: "http://kuzzle:7512",
						Routes: Routes{
							Login:          "/_login/local",
							GetCurrentUser: "/_me",
						},
						AllowedUsers: []string{"test"},
					},
					CustomRealm: "Use a valid Kuzzle user to authenticate",
				},
			},
			expNextCall:   true,
			expStatusCode: 200,
			basicAuthMock: struct {
				enabled  bool
				username string
				password string
			}{
				enabled:  true,
				username: "test",
				password: "test",
			},
			loginMock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_login/local",
				response:   json.RawMessage(`{"result":{"jwt": "myToken"}}`),
			},
			getCurrentUserMock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_me",
				response:   json.RawMessage(`{"result":{"_id": "test"}}`),
			},
		},
		{
			name: "Success without allowedUsers configured",
			fields: fields{
				name: "kuzzle-auth",
				config: &Config{
					Kuzzle: Kuzzle{
						URL: "http://kuzzle:7512",
						Routes: Routes{
							Login:          "/_login/local",
							GetCurrentUser: "/_me",
						},
					},
					CustomRealm: "Use a valid Kuzzle user to authenticate",
				},
			},
			expNextCall:   true,
			expStatusCode: 200,
			basicAuthMock: struct {
				enabled  bool
				username string
				password string
			}{
				enabled:  true,
				username: "test",
				password: "test",
			},
			loginMock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_login/local",
				response:   json.RawMessage(`{"result":{"jwt": "myToken"}}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.loginMock.enabled {
				defer gock.Off()
				gock.
					New(tt.loginMock.url).
					Post(tt.loginMock.route).
					Reply(tt.loginMock.statusCode).
					JSON(tt.loginMock.response)
			}

			if tt.getCurrentUserMock.enabled {
				defer gock.Off()
				gock.
					New(tt.getCurrentUserMock.url).
					Get(tt.getCurrentUserMock.route).
					Reply(tt.getCurrentUserMock.statusCode).
					JSON(tt.getCurrentUserMock.response)
			}

			nextCall := false
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				nextCall = true
			})

			ka := &KuzzleAuth{
				next:   next,
				name:   tt.fields.name,
				config: tt.fields.config,
			}

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			if tt.basicAuthMock.enabled {
				req.SetBasicAuth(tt.basicAuthMock.username, tt.basicAuthMock.password)
			}

			ka.ServeHTTP(recorder, req)

			if nextCall != tt.expNextCall {
				t.Errorf("next handler should not be called")
			}

			if recorder.Result().StatusCode != tt.expStatusCode {
				t.Errorf("got status code %d, want %d", recorder.Code, tt.expStatusCode)
			}
			ka.ServeHTTP(recorder, req)
		})
	}
}
