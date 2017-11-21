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
