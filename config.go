package traefik_kuzzle_auth

// Config for the plugin configuration.
type Config struct {
	Kuzzle      Kuzzle `yaml:"kuzzle"`      // Kuzzle remote server configuration
	CustomRealm string `yaml:"customRealm"` // CustomRealm can be used to personalize Basic Auth window message
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	c := &Config{}
	return c.addMissingFields()
}

func (c *Config) addMissingFields() *Config {
	if c.CustomRealm == "" {
		c.CustomRealm = "Use a valid Kuzzle user to authenticate"
	}

	if c.Kuzzle.Routes.Login == "" {
		c.Kuzzle.Routes.Login = "/_login/local"
	}

	if c.Kuzzle.Routes.GetCurrentUser == "" {
		c.Kuzzle.Routes.GetCurrentUser = "/_me"
	}

	return c
}
