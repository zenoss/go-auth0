package mgmt

import (
	"github.com/zenoss/go-auth0/auth0/http"
)

// ManagementService is a gateway to Auth0 Management services
type ManagementService struct {
	*http.Client
	Users *UsersService
}

// New creates a new ManagementService, backed by client
func New(site string, client *http.Client) *ManagementService {
	mgmt := &ManagementService{
		Client: &http.Client{
			Doer: client,
			Site: site,
		},
	}
	mgmt.Users = &UsersService{
		c: mgmt.Client,
	}
	return mgmt
}
