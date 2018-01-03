package mgmt

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

// ManagementService is a gateway to Auth0 Management services
type ManagementService struct {
	Client
	Site  string
	Users *UsersService
}

// New creates a new ManagementService, backed by client
func New(site string, client Client) *ManagementService {
	mgmt := &ManagementService{
		Client: client,
		Site:   site,
	}
	mgmt.Users = &UsersService{
		c: mgmt,
	}
	return mgmt
}

// Do processes a request and unmarshals the response body into respBody
func (mgmt *ManagementService) Do(req *http.Request, respBody interface{}) error {
	return mgmt.Client.Do(req, respBody)
}

// Get performs a get to the endpoint of the site associated with the client
func (mgmt *ManagementService) Get(endpoint string, respBody interface{}) error {
	req, err := http.NewRequest("GET", mgmt.Site+endpoint, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return mgmt.Do(req, respBody)
}

// Post performs a post to the endpoint of the site associated with the client
func (mgmt *ManagementService) Post(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("POST", mgmt.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return mgmt.Do(req, respBody)
}

// Put performs a put to the endpoint of the site associated with the client
func (mgmt *ManagementService) Put(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("PUT", mgmt.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return mgmt.Do(req, respBody)
}

// Patch performs a patch to the endpoint of the site associated with the client
func (mgmt *ManagementService) Patch(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("PATCH", mgmt.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return mgmt.Do(req, respBody)
}

// Delete performs a delete to the endpoint of the site associated with the client
func (mgmt *ManagementService) Delete(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("DELETE", mgmt.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return mgmt.Do(req, respBody)
}
