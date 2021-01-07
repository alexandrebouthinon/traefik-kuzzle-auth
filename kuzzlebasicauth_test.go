package kuzzlebasicauth

import (
	"context"
	"net/http"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func TestConfig_Check(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid HTTP Configuration",
			fields: fields{
				Kuzzle: struct {
					Host string `json:"host,omitempty"`
					Port uint   `json:"port,omitempty"`
					Ssl  bool   `json:"ssl,omitempty"`
					URL  string
				}{
					Host: "localhost",
					Port: 7512,
					Ssl:  false,
				},
				BasicAuth: struct {
					User     string `json:"user,omitempty"`
					Password string `json:"password,omitempty"`
				}{
					User:     "user",
					Password: "password",
				},
			},
			wantErr: false,
		},
		{
			name: "Kuzzle Host empty string",
			fields: fields{
				Kuzzle: struct {
					Host string `json:"host,omitempty"`
					Port uint   `json:"port,omitempty"`
					Ssl  bool   `json:"ssl,omitempty"`
					URL  string
				}{
					Host: "",
					Port: 7512,
					Ssl:  false,
				},
				BasicAuth: struct {
					User     string `json:"user,omitempty"`
					Password string `json:"password,omitempty"`
				}{
					User:     "user",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "Basic Auth user empty string",
			fields: fields{
				Kuzzle: struct {
					Host string `json:"host,omitempty"`
					Port uint   `json:"port,omitempty"`
					Ssl  bool   `json:"ssl,omitempty"`
					URL  string
				}{
					Host: "localhost",
					Port: 7512,
					Ssl:  false,
				},
				BasicAuth: struct {
					User     string `json:"user,omitempty"`
					Password string `json:"password,omitempty"`
				}{
					User:     "",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "Basic Auth password empty string",
			fields: fields{
				Kuzzle: struct {
					Host string `json:"host,omitempty"`
					Port uint   `json:"port,omitempty"`
					Ssl  bool   `json:"ssl,omitempty"`
					URL  string
				}{
					Host: "localhost",
					Port: 7512,
					Ssl:  false,
				},
				BasicAuth: struct {
					User     string `json:"user,omitempty"`
					Password string `json:"password,omitempty"`
				}{
					User:     "user",
					Password: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Kuzzle:    tt.fields.Kuzzle,
				BasicAuth: tt.fields.BasicAuth,
			}

			if err := c.check(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_PingKuzzle(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
		useMock bool
		mock    struct {
			enabled bool
			status  int
		}
	}{
		{
			name: "Kuzzle HTTP ping successful",
			fields: fields{
				Kuzzle: struct {
					Host string `json:"host,omitempty"`
					Port uint   `json:"port,omitempty"`
					Ssl  bool   `json:"ssl,omitempty"`
					URL  string
				}{
					Host: "localhost",
					Port: 7512,
					Ssl:  false,
				},
				BasicAuth: struct {
					User     string `json:"user,omitempty"`
					Password string `json:"password,omitempty"`
				}{
					User:     "user",
					Password: "password",
				},
			},
			wantErr: false,
			want:    "http://localhost:7512",
			mock: struct {
				enabled bool
				status  int
			}{enabled: true, status: 200},
		},
		{
			name: "Kuzzle HTTPS ping successful",
			fields: fields{
				Kuzzle: struct {
					Host string `json:"host,omitempty"`
					Port uint   `json:"port,omitempty"`
					Ssl  bool   `json:"ssl,omitempty"`
					URL  string
				}{
					Host: "localhost",
					Port: 7512,
					Ssl:  true,
				},
				BasicAuth: struct {
					User     string `json:"user,omitempty"`
					Password string `json:"password,omitempty"`
				}{
					User:     "user",
					Password: "password",
				},
			},
			wantErr: false,
			want:    "https://localhost:7512",
			mock: struct {
				enabled bool
				status  int
			}{enabled: true, status: 200},
		},
		{
			name: "Kuzzle HTTPS ping failure",
			fields: fields{
				Kuzzle: struct {
					Host string `json:"host,omitempty"`
					Port uint   `json:"port,omitempty"`
					Ssl  bool   `json:"ssl,omitempty"`
					URL  string
				}{
					Host: "nowhere",
					Port: 7512,
					Ssl:  true,
				},
				BasicAuth: struct {
					User     string `json:"user,omitempty"`
					Password string `json:"password,omitempty"`
				}{
					User:     "user",
					Password: "password",
				},
			},
			wantErr: true,
			want:    "https://nowhere:7512",
			mock: struct {
				enabled bool
				status  int
			}{enabled: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock.enabled {
				defer gock.Off()
				gock.New(tt.want).Get("/_public").Reply(tt.mock.status)
			}

			c := &Config{
				Kuzzle:    tt.fields.Kuzzle,
				BasicAuth: tt.fields.BasicAuth,
			}

			got, err := pingKuzzle(c)

			if (err != nil) != tt.wantErr {
				t.Errorf("Config.PingKuzzle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Config.PingKuzzle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKuzzleBasicAuth_loginToKuzzleIsSuccessful(t *testing.T) {
	type fields struct {
		next   http.Handler
		name   string
		config *Config
	}
	type args struct {
		user     string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
		mock    struct {
			enabled bool
			status  int
		}
	}{
		{
			name: "Kuzzle login successful",
			args: struct {
				user     string
				password string
			}{
				user:     "user",
				password: "password",
			},
			fields: fields{
				config: &Config{
					Kuzzle: struct {
						Host string `json:"host,omitempty"`
						Port uint   `json:"port,omitempty"`
						Ssl  bool   `json:"ssl,omitempty"`
						URL  string
					}{
						Host: "localhost",
						Port: 7512,
						Ssl:  true,
						URL:  "https://localhost:7512",
					},
					BasicAuth: struct {
						User     string `json:"user,omitempty"`
						Password string `json:"password,omitempty"`
					}{
						User:     "user",
						Password: "password",
					},
				},
			},
			wantErr: false,
			want:    true,
			mock: struct {
				enabled bool
				status  int
			}{enabled: true, status: 200},
		},
		{
			name: "Kuzzle login failure: unreachable",
			args: struct {
				user     string
				password string
			}{
				user:     "user",
				password: "password",
			},
			fields: fields{
				config: &Config{
					Kuzzle: struct {
						Host string `json:"host,omitempty"`
						Port uint   `json:"port,omitempty"`
						Ssl  bool   `json:"ssl,omitempty"`
						URL  string
					}{
						Host: "localhost",
						Port: 7512,
						Ssl:  true,
						URL:  "https://localhost:7512",
					},
					BasicAuth: struct {
						User     string `json:"user,omitempty"`
						Password string `json:"password,omitempty"`
					}{
						User:     "user",
						Password: "password",
					},
				},
			},
			wantErr: true,
			want:    false,
			mock: struct {
				enabled bool
				status  int
			}{enabled: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock.enabled {
				defer gock.Off()
				gock.New(tt.fields.config.Kuzzle.URL).Post("/_login/local").Reply(tt.mock.status)
			}

			k := &KuzzleBasicAuth{
				next:   tt.fields.next,
				name:   tt.fields.name,
				config: tt.fields.config,
			}
			got, err := k.loginToKuzzleIsSuccessful(tt.args.user, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("KuzzleBasicAuth.loginToKuzzleIsSuccessful() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("KuzzleBasicAuth.loginToKuzzleIsSuccessful() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		wantErr bool
		mock    struct {
			enabled bool
			status  int
		}
	}{
		{
			name: "Plugin creation successful",
			args: args{
				ctx:  nil,
				next: nil,
				name: "Plugin",
				config: &Config{
					Kuzzle: struct {
						Host string `json:"host,omitempty"`
						Port uint   `json:"port,omitempty"`
						Ssl  bool   `json:"ssl,omitempty"`
						URL  string
					}{
						Host: "localhost",
						Port: 7512,
						Ssl:  false,
					},
					BasicAuth: struct {
						User     string `json:"user,omitempty"`
						Password string `json:"password,omitempty"`
					}{
						User:     "user",
						Password: "password",
					},
				},
			},
			wantErr: false,
			mock: struct {
				enabled bool
				status  int
			}{
				enabled: true,
				status:  200,
			},
		},
		{
			name: "Configuration check failure",
			args: args{
				ctx:  nil,
				next: nil,
				name: "Plugin",
				config: &Config{
					Kuzzle: struct {
						Host string `json:"host,omitempty"`
						Port uint   `json:"port,omitempty"`
						Ssl  bool   `json:"ssl,omitempty"`
						URL  string
					}{
						Host: "",
						Port: 7512,
						Ssl:  false,
					},
					BasicAuth: struct {
						User     string `json:"user,omitempty"`
						Password string `json:"password,omitempty"`
					}{
						User:     "user",
						Password: "password",
					},
				},
			},
			wantErr: true,
			mock: struct {
				enabled bool
				status  int
			}{
				enabled: true,
				status:  200,
			},
		},
		{
			name: "Kuzzle ping failure",
			args: args{
				ctx:  nil,
				next: nil,
				name: "Plugin",
				config: &Config{
					Kuzzle: struct {
						Host string `json:"host,omitempty"`
						Port uint   `json:"port,omitempty"`
						Ssl  bool   `json:"ssl,omitempty"`
						URL  string
					}{
						Host: "localhost",
						Port: 7512,
						Ssl:  false,
					},
					BasicAuth: struct {
						User     string `json:"user,omitempty"`
						Password string `json:"password,omitempty"`
					}{
						User:     "user",
						Password: "password",
					},
				},
			},
			wantErr: true,
			mock: struct {
				enabled bool
				status  int
			}{enabled: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock.enabled {
				defer gock.Off()
				gock.New(tt.args.config.Kuzzle.URL).Get("/_publicApi").Reply(tt.mock.status)
			}

			_, err := New(tt.args.ctx, tt.args.next, tt.args.config, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
