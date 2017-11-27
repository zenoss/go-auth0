package auth0

import (
	"fmt"
	"net/http"
	"time"
)

// Auth0 is a client used to make requests to Auth0 APIs
type Auth0 struct {
	c    *http.Client
	site string
}

// NewClient creates an Auth0 client
func NewClient(tenant string) *Auth0 {
	return &Auth0{
		c: &http.Client{
			Timeout: time.Second * 5,
		},
		site: fmt.Sprintf("https://%s.auth0.com", tenant),
	}
}

func (auth *Auth0) do(req *http.Request) (*http.Response, error) {
	// POSTs are application/json to this api
	if req.Method == "POST" {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := auth.c.Do(req)
	return resp, err
}
