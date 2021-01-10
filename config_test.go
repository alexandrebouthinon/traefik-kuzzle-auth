package traefik_kuzzle_auth

import (
	"reflect"
	"testing"
)

func TestConfig_addMissingFields(t *testing.T) {
	type fields struct {
		Kuzzle      Kuzzle
		CustomRealm string
	}
	tests := []struct {
		name   string
		fields fields
		want   *Config
	}{
		{
			name: "Complete valid configuration",
			fields: fields{
				Kuzzle: Kuzzle{
					URL: "http://kuzzle:7512",
					Routes: Routes{
						Login:          "/_login/local",
						GetCurrentUser: "/_me",
					},
					AllowedUsers: []string{"admin"},
				},
				CustomRealm: "Use valid user to authenticate",
			},
			want: &Config{
				Kuzzle: Kuzzle{
					URL: "http://kuzzle:7512",
					Routes: Routes{
						Login:          "/_login/local",
						GetCurrentUser: "/_me",
					},
					AllowedUsers: []string{"admin"},
				},
				CustomRealm: "Use valid user to authenticate",
			},
		},
		{
			name: "Missing all optional fields",
			fields: fields{
				Kuzzle: Kuzzle{
					URL:    "http://kuzzle:7512",
					Routes: Routes{},
				},
			},
			want: &Config{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Kuzzle:      tt.fields.Kuzzle,
				CustomRealm: tt.fields.CustomRealm,
			}
			if got := c.addMissingFields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.addMissingFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
