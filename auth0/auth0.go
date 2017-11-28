package auth0

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spopezen/go-auth0/auth0/authz"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
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

// Config is used to create a client that uses the OAuth2 flow
type Config struct {
	Tenant           string
	ClientID         string
	ClientSecret     string
	RedirectURI      string
	AuthorizationAPI string
	Scopes           []string
}

// Auth0 is a client used to make requests to Auth0 APIs
type Auth0 struct {
	c     *http.Client
	Site  string
	Token *TokenService
	Authz *authz.AuthorizationService
}

// GrantFunc gets an Authorization Grant
type GrantFunc func(url string) (string, error)

// PromptGrant uses stdin/out to get an authorization grant
func PromptGrant(url string) (string, error) {
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return "", errors.Wrap(err, "Unable to get code from input")
	}
	return code, nil
}

// ClientFromGrant creates a client that follows the standard "3-legged" OAuth2 flow
func (conf *Config) ClientFromGrant(getGrant GrantFunc) (*Auth0, error) {
	ctx := context.Background()
	cfg := &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s.auth0.com/authorize", conf.Tenant),
			TokenURL: fmt.Sprintf("https://%s.auth0.com/oauth/token", conf.Tenant),
		},
		Scopes: conf.Scopes,
	}
	url := cfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
	code, err := getGrant(url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get authorization grant")
	}
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to exchange authorization grant for token")
	}
	c := &Auth0{
		c:    cfg.Client(ctx, token),
		Site: fmt.Sprintf("https://%s.auth0.com", conf.Tenant),
	}
	c.Token = &TokenService{
		Client: c,
	}
	c.Authz = authz.New(conf.AuthorizationAPI, c)
	return c, nil
}

// ClientFromCredentials creates a client that follows the "2-legged"
// Client Credientials OAuth2 flow
func (conf *Config) ClientFromCredentials(API string) (*Auth0, error) {
	ctx := context.Background()
	cfg := &clientcredentials.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		TokenURL:     fmt.Sprintf("https://%s.auth0.com/oauth/token", conf.Tenant),
		Scopes:       conf.Scopes,
		EndpointParams: url.Values{
			"audience": []string{API},
		},
	}
	c := &Auth0{
		c:    cfg.Client(ctx),
		Site: fmt.Sprintf("https://%s.auth0.com", conf.Tenant),
	}
	c.Token = &TokenService{
		Client: c,
	}
	c.Authz = authz.New(conf.AuthorizationAPI, c)
	return c, nil
}

// Do processes a request and unmarshals the response body into respBody
func (auth *Auth0) Do(req *http.Request, respBody interface{}) error {
	// POSTs are application/json to this api
	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}
	// fmt.Printf("Request:\n%+v\n\n", req)
	resp, err := auth.c.Do(req)
	if err != nil {
		return errors.Wrap(err, "Cannot complete request")
	}
	fmt.Printf("Response Status: %s\n", resp.Status)
	if resp.ContentLength == 0 {
		fmt.Println("No body in response")
		return nil
	}
	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Cannot read response body")
	}
	fmt.Printf("Response Body:\n%s\n\n", bodyData)
	err = json.Unmarshal(bodyData, respBody)
	if err != nil {
		return errors.Wrap(err, "Cannot unmarshal response")
	}
	return nil
}

// Get performs a get to the endpoint of the site associated with the client
func (auth *Auth0) Get(endpoint string, respBody interface{}) error {
	req, err := http.NewRequest("GET", auth.Site+endpoint, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return auth.Do(req, respBody)
}

// Post performs a post to the endpoint of the site associated with the client
func (auth *Auth0) Post(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("POST", auth.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return auth.Do(req, respBody)
}

// Put performs a put to the endpoint of the site associated with the client
func (auth *Auth0) Put(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("PUT", auth.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return auth.Do(req, respBody)
}

// Patch performs a patch to the endpoint of the site associated with the client
func (auth *Auth0) Patch(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("PATCH", auth.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return auth.Do(req, respBody)
}

// Delete performs a delete to the endpoint of the site associated with the client
func (auth *Auth0) Delete(endpoint string, body interface{}, respBody interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("DELETE", auth.Site+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	return auth.Do(req, respBody)
}
