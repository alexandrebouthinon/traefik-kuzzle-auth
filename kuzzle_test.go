package traefik_kuzzle_auth

import (
	"encoding/json"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func TestKuzzle_login(t *testing.T) {
	type fields struct {
		URL          string
		Routes       Routes
		AllowedUsers []string
		JWT          string
	}
	type args struct {
		user     string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    Mock
	}{
		{
			name: "Success",
			fields: fields{
				URL: "http://kuzzle:7512",
				Routes: Routes{
					Login: "/_login/local",
				},
			},
			wantErr: false,
			mock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_login/local",
				response:   json.RawMessage(`{"result":{"jwt": "myToken"}}`),
			},
			args: args{
				user:     "user",
				password: "pass",
			},
		},
		{
			name: "Fail sending login request",
			fields: fields{
				URL: "http://kuzzle:7512",
				Routes: Routes{
					Login: "/_login/local",
				},
			},
			wantErr: true,
			args: args{
				user:     "user",
				password: "pass",
			},
		},
		{
			name: "Bad credentials or unauthorized route",
			fields: fields{
				URL: "http://kuzzle:7512",
				Routes: Routes{
					Login: "/_login/local",
				},
			},
			wantErr: true,
			mock: Mock{
				enabled:    true,
				statusCode: 401,
				url:        "http://kuzzle:7512",
				route:      "/_login/local",
				response:   json.RawMessage(`{"result":{"jwt": "myToken"}}`),
			},
			args: args{
				user:     "user",
				password: "pass",
			},
		},
		{
			name: "Wrong JSON response format",
			fields: fields{
				URL: "http://kuzzle:7512",
				Routes: Routes{
					Login: "/_login/local",
				},
			},
			wantErr: true,
			mock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_login/local",
				response:   `{"result"}`,
			},
			args: args{
				user:     "user",
				password: "pass",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock.enabled {
				defer gock.Off()
				gock.
					New(tt.mock.url).
					Post(tt.mock.route).
					Reply(tt.mock.statusCode).
					JSON(tt.mock.response)
			}
			k := &Kuzzle{
				URL:          tt.fields.URL,
				Routes:       tt.fields.Routes,
				AllowedUsers: tt.fields.AllowedUsers,
				JWT:          tt.fields.JWT,
			}
			if err := k.login(tt.args.user, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("Kuzzle.login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKuzzle_checkUser(t *testing.T) {
	type fields struct {
		URL          string
		Routes       Routes
		AllowedUsers []string
		JWT          string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		mock    Mock
	}{
		{
			name: "Success",
			fields: fields{
				URL: "http://kuzzle:7512",
				Routes: Routes{
					GetCurrentUser: "/_me",
				},
				AllowedUsers: []string{"admin"},
				JWT:          "validJWT",
			},
			wantErr: false,
			mock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_me",
				response:   json.RawMessage(`{"result":{"_id": "admin"}}`),
			},
		},
		{
			name: "Unreachable Kuzzle server",
			fields: fields{
				URL: "http://kuzzle:7512",
				Routes: Routes{
					GetCurrentUser: "/_me",
				},
				AllowedUsers: []string{"admin"},
				JWT:          "validJWT",
			},
			wantErr: true,
		},
		{
			name: "Wrong JSON response format",
			fields: fields{
				URL: "http://kuzzle:7512",
				Routes: Routes{
					GetCurrentUser: "/_me",
				},
				AllowedUsers: []string{"admin"},
				JWT:          "validJWT",
			},
			wantErr: true,
			mock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_me",
				response:   `{"result"}`,
			},
		},
		{
			name: "Not part of allowed users",
			fields: fields{
				URL: "http://kuzzle:7512",
				Routes: Routes{
					GetCurrentUser: "/_me",
				},
				AllowedUsers: []string{"foo", "bar"},
				JWT:          "validJWT",
			},
			wantErr: true,
			mock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_me",
				response:   json.RawMessage(`{"result":{"_id": "admin"}}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock.enabled {
				defer gock.Off()
				gock.
					New(tt.mock.url).
					Get(tt.mock.route).
					Reply(tt.mock.statusCode).
					JSON(tt.mock.response)
			}

			k := &Kuzzle{
				URL:          tt.fields.URL,
				Routes:       tt.fields.Routes,
				AllowedUsers: tt.fields.AllowedUsers,
				JWT:          tt.fields.JWT,
			}
			if err := k.checkUser(); (err != nil) != tt.wantErr {
				t.Errorf("Kuzzle.checkUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
