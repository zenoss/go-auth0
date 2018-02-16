package auth0

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/zenoss/go-auth0/auth0/authz"
	"github.com/zenoss/go-auth0/auth0/http"
	"github.com/zenoss/go-auth0/auth0/mgmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Config is used to create a client that uses the OAuth2 flow
type Config struct {
	Tenant           string
	ClientID         string
	ClientSecret     string
	RedirectURI      string
	AuthorizationURL string
	Scopes           []string
}

// Auth0 is a client used to make requests to Auth0 APIs
type Auth0 struct {
	*http.Client
	Token *TokenService
	Authz *authz.AuthorizationService
	Mgmt  *mgmt.ManagementService
}

// GrantFunc gets an Authorization Grant
type GrantFunc func(URL string) (string, error)

// PromptGrant uses stdin/out to get an authorization grant
func PromptGrant(URL string) (string, error) {
	fmt.Printf("Visit the URL for the auth dialog: %v\n", URL)
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
	URL := cfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
	code, err := getGrant(URL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get authorization grant")
	}
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to exchange authorization grant for token")
	}
	c := &Auth0{
		Client: &http.Client{
			Doer: &http.RootClient{
				Client: cfg.Client(ctx, token),
			},
			Site: fmt.Sprintf("https://%s.auth0.com", conf.Tenant),
		},
	}
	c.Token = &TokenService{
		Client: c.Client,
	}
	c.Mgmt = mgmt.New(c.Site, c.Client)
	if conf.AuthorizationURL != "" {
		c.Authz = authz.New(conf.AuthorizationURL, c.Client)
	}
	return c, nil
}

// ClientFromCredentials creates a client that follows the "2-legged"
// Client Credientials OAuth2 flow
func (conf *Config) ClientFromCredentials(APIs []string) (*Auth0, error) {
	ctx := context.Background()
	cfg := &clientcredentials.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		TokenURL:     fmt.Sprintf("https://%s.auth0.com/oauth/token", conf.Tenant),
		Scopes:       conf.Scopes,
		EndpointParams: url.Values{
			"audience": APIs,
		},
	}
	c := &Auth0{
		Client: &http.Client{
			Doer: &http.RootClient{
				Client: cfg.Client(ctx),
			},
			Site: fmt.Sprintf("https://%s.auth0.com", conf.Tenant),
		},
	}
	c.Token = &TokenService{
		Client: c.Client,
	}
	c.Mgmt = mgmt.New(c.Site, c.Client)
	if conf.AuthorizationURL != "" {
		c.Authz = authz.New(conf.AuthorizationURL, c.Client)
	}
	return c, nil
}
