package traefik_kuzzle_auth

// Mock configuration used to set up gock HTTP mocking
type Mock struct {
	enabled    bool
	statusCode int
	url        string
	route      string
	response   interface{}
}
