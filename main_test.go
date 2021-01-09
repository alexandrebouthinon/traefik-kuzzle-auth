package traefik_kuzzle_auth

import (
	"context"
	"net/http"
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
						URL: "http://kuzzle:7512",
						Routes: Routes{
							Ping: "/_publicApi",
						},
					},
				},
			},
			want: &KuzzleAuth{
				config: &Config{
					Kuzzle: Kuzzle{
						URL: "http://kuzzle:7512",
						Routes: Routes{
							Ping: "/_publicApi",
						},
					},
				},
			},
			wantErr: false,
			mock: Mock{
				enabled:    true,
				statusCode: 200,
				url:        "http://kuzzle:7512",
				route:      "/_publicApi",
			},
		},
		{
			name: "Kuzzle server ping failure",
			args: args{
				config: &Config{
					Kuzzle: Kuzzle{
						URL: "http://kuzzle:7512",
						Routes: Routes{
							Ping: "/_publicApi",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
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
