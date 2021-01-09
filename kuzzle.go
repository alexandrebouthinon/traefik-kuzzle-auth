package traefik_kuzzle_auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Routes used to request Kuzzle, can be customized
type Routes struct {
	// Ping route used to test configured Kuzzle server reachability.
	// The specified route must return 200 HTTP status code when called by anonymous user.
	// Default is '/_publicApi' (see: https://docs.kuzzle.io/core/2/api/controllers/server/public-api/).
	// You would like to modify this route if you performed security adjustement on your Kuzzle Server
	Ping string `yaml:"ping,omitempty"`
	// Login route used to log in to Kuzzle using Auth Basic user/pass.
	// The specified route must return 200 HTTP status code and a valid JWT when called by anonymous user.
	// Default is '/_login/local' (see: https://docs.kuzzle.io/core/2/api/controllers/auth/login/)
	// Login route using 'local' strtategy (see: https://docs.kuzzle.io/core/2/guides/main-concepts/authentication/#local-strategy)
	// It must accept JSON body containing 'username' and 'password' string fields, for example:
	// 	{
	// 		"username": "myUser",
	// 		"password": "myV3rys3cretP4ssw0rd"
	// 	}
	// You would like to update this route if you do not use 'local' strategy on your Kuzzle server
	Login string `yaml:"login,omitempty"`
	// GetCurrentUser route used to get logged in user KUID.
	// Default is '/_me' but this a Kuzzle v2 route only so you would like update it if you still use Kuzzle v1
	// (see: https://docs.kuzzle.io/core/2/api/controllers/auth/get-current-user/).
	GetCurrentUser string `yaml:"getCurrentUser,omitempty"`
}

// Kuzzle info
type Kuzzle struct {
	// URL use by the plugin to reach Kuzzle server.
	// NOTE: Only HTTP(s) protocol is supported
	// Examples:
	//  - HTTP: http://localhost:7512
	//	- HTTPS: https://localhost:7512
	URL    string `yaml:"url"`
	Routes Routes `yaml:"routes,omitempty"`
	// AllowedUsers contain users KUID allowed to connect using this plugin.
	// It is empty by default so every user registered on your Kuzzle server can use this plugin.
	// More about users KUID at https://docs.kuzzle.io/core/2/guides/main-concepts/authentication/#kuzzle-user-identifier-kuid
	// NOTE: The user you used to log in need to be able to call `auth:getCurrentUser` Kuzzle API route
	AllowedUsers []string `yaml:"allowedUsers,omitempty"`
	JWT          string
}

func (k *Kuzzle) ping() error {
	url := fmt.Sprintf("%s%s", k.URL, k.Routes.Ping)
	_, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("Ping request send to %s failed: %v", url, err)
	}

	return nil
}

func (k *Kuzzle) login(user string, password string) error {
	reqBody, _ := json.Marshal(map[string]string{
		"username": user,
		"password": password,
	})

	url := fmt.Sprintf("%s%s", k.URL, k.Routes.Login)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	if err != nil {
		return fmt.Errorf("Authentication request send to %s failed: %v", url, err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Authentication request send to %s failed: status code %d", url, resp.StatusCode)
	}

	var jsonBody map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return err
	}

	k.JWT = jsonBody["result"].(map[string]interface{})["jwt"].(string)

	return nil
}

func (k *Kuzzle) checkUser() error {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s", k.URL, k.Routes.GetCurrentUser)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", k.JWT))
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	var jsonBody map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return err
	}

	kuid := jsonBody["result"].(map[string]interface{})["_id"].(string)
	for _, id := range k.AllowedUsers {
		if kuid == id {
			return nil
		}
	}

	return fmt.Errorf("User %s do not be part of allowed users: %v", kuid, k.AllowedUsers)
}
