package authz

import (
	"github.com/zenoss/go-auth0/auth0/http"
)

// AuthorizationService is a gateway to Auth0 Authorization Extension services
type AuthorizationService struct {
	*http.Client
	Groups      *GroupsService
	Permissions *PermissionsService
	Roles       *RolesService
	Users       *UsersService
}

// New creates a new AuthorizationService, backed by client, whose
// Authorization extension lives at site
func New(client *http.Client) *AuthorizationService {
	authz := &AuthorizationService{
		Client: client,
	}
	authz.Groups = &GroupsService{
		c: authz.Client,
	}
	authz.Permissions = &PermissionsService{
		c: authz.Client,
	}
	authz.Roles = &RolesService{
		c: authz.Client,
	}
	authz.Users = &UsersService{
		c: authz.Client,
	}
	return authz
}
