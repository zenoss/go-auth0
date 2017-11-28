package authz

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Client can perform https requests
type Client interface {
	Do(req *http.Request, respBody interface{}) error
	Get(endpoint string, respBody interface{}) error
	Post(endpoint string, body interface{}, respBody interface{}) error
	Put(endpoint string, body interface{}, respBody interface{}) error
	Patch(endpoint string, body interface{}, respBody interface{}) error
	Delete(endpoint string, body interface{}, respBody interface{}) error
}

// AuthorizationService is a gateway to Auth0 Authorization Extension services
type AuthorizationService struct {
	Client
	Site        string
	Groups      *GroupsService
	Permissions *PermissionsService
	Roles       *RolesService
	Users       *UsersService
}

// New creates a new AuthorizationService, backed by client, whose
// Authorization extension lives at site
func New(site string, client Client) *AuthorizationService {
	authz := &AuthorizationService{
		Client: client,
		Site:   site,
	}
	authz.Groups = &GroupsService{
		c: authz,
	}
	authz.Permissions = &PermissionsService{
		Client: authz,
	}
	authz.Roles = &RolesService{
		Client: authz,
	}
	authz.Users = &UsersService{
		Client: authz,
	}
	return authz
}

// Do processes a request and unmarshals the response body into respBody
func (authz *AuthorizationService) Do(req *http.Request, respBody interface{}) error {
	return authz.Client.Do(req, respBody)
}

// Get performs a get to the endpoint of the site associated with the client
func (authz *AuthorizationService) Get(endpoint string, respBody interface{}) error {
	req, err := http.NewRequest("GET", authz.Site+endpoint, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return authz.Do(req, respBody)
}

// Post performs a post to the endpoint of the site associated with the client
func (authz *AuthorizationService) Post(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("POST", authz.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return authz.Do(req, respBody)
}

// Put performs a put to the endpoint of the site associated with the client
func (authz *AuthorizationService) Put(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("PUT", authz.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return authz.Do(req, respBody)
}

// Patch performs a patch to the endpoint of the site associated with the client
func (authz *AuthorizationService) Patch(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("PATCH", authz.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return authz.Do(req, respBody)
}

// Delete performs a delete to the endpoint of the site associated with the client
func (authz *AuthorizationService) Delete(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("DELETE", authz.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return authz.Do(req, respBody)
}
