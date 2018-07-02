package auth0

import (
	"context"
	"fmt"
	gohttp "net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/zenoss/go-auth0/auth0/authz"
	"github.com/zenoss/go-auth0/auth0/http"
	"github.com/zenoss/go-auth0/auth0/mgmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// TokenClient is a client to the token endpoint
func TokenClient(domain string) *TokenService {
	return &TokenService{
		&http.Client{
			Doer: &http.RootClient{
				Client: &gohttp.Client{},
			},
			API: fmt.Sprintf("https://%s", domain),
		},
	}
}

// API represents an api in Auth0
type API struct {
	URL          string
	Audience     []string
	ClientID     string
	ClientSecret string
	Scopes       []string
}

func clientCredentialsConfig(ctx context.Context, domain string, api API) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     api.ClientID,
		ClientSecret: api.ClientSecret,
		TokenURL:     fmt.Sprintf("https://%s/oauth/token", domain),
		Scopes:       api.Scopes,
		EndpointParams: url.Values{
			"audience": api.Audience,
		},
	}
}

// ClientFromCredentials follows the returns a go http Client authorized for the given API
func ClientFromCredentials(domain string, api API) *gohttp.Client {
	ctx := context.Background()
	return clientCredentialsConfig(ctx, domain, api).Client(ctx)
}

// MgmtClientFromCredentials follows the returns a client authorized for the given API
func MgmtClientFromCredentials(domain string, api API) *mgmt.ManagementService {
	ctx := context.Background()
	cfg := clientCredentialsConfig(ctx, domain, api)
	return mgmt.New(
		&http.Client{
			Doer: &http.RootClient{
				Client: cfg.Client(ctx),
			},
			API: api.URL,
		},
	)
}

// AuthzClientFromCredentials follows the returns a client authorized for the given API
func AuthzClientFromCredentials(domain string, api API) *authz.AuthorizationService {
	ctx := context.Background()
	cfg := clientCredentialsConfig(ctx, domain, api)
	return authz.New(
		&http.Client{
			Doer: &http.RootClient{
				Client: cfg.Client(ctx),
			},
			API: api.URL,
		},
	)
}

// GrantFunc is a function that gets an Authorization Grant code
type GrantFunc func(URL string) (string, error)

// PromptGrant uses stdin/out to get an authorization grant
func PromptGrant(URL string) (string, error) {
	fmt.Printf("Visit the URL for the auth dialog, then enter the code: %v\n", URL)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return "", errors.Wrap(err, "Unable to get code from input")
	}

	return code, nil
}

func grantConfig(ctx context.Context, domain string, api API) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     api.ClientID,
		ClientSecret: api.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s/authorize", domain),
			TokenURL: fmt.Sprintf("https://%s/oauth/token", domain),
		},
		Scopes: api.Scopes,
	}
}

// ClientFromGrant follows the 3-legged OAuth2 flow to get an authorized client for the given API
func ClientFromGrant(domain string, api API, getGrant GrantFunc) (*gohttp.Client, error) {
	ctx := context.Background()
	cfg := grantConfig(ctx, domain, api)
	URL := cfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
	code, err := getGrant(URL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get authorization grant")
	}
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to exchange authorization grant for token")
	}
	return cfg.Client(ctx, token), nil
}

// MgmtClientFromGrant follows the 3-legged OAuth2 flow to get an authorized client for the given API
func MgmtClientFromGrant(domain string, api API, getGrant GrantFunc) (*mgmt.ManagementService, error) {
	ctx := context.Background()
	cfg := grantConfig(ctx, domain, api)
	URL := cfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
	code, err := getGrant(URL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get authorization grant")
	}
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to exchange authorization grant for token")
	}
	return mgmt.New(
		&http.Client{
			Doer: &http.RootClient{
				Client: cfg.Client(ctx, token),
			},
			API: api.URL,
		},
	), nil
}

// AuthzClientFromGrant follows the 3-legged OAuth2 flow to get an authorized client for the given API
func AuthzClientFromGrant(domain string, api API, getGrant GrantFunc) (*authz.AuthorizationService, error) {
	ctx := context.Background()
	cfg := grantConfig(ctx, domain, api)
	URL := cfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
	code, err := getGrant(URL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get authorization grant")
	}
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to exchange authorization grant for token")
	}
	return authz.New(
		&http.Client{
			Doer: &http.RootClient{
				Client: cfg.Client(ctx, token),
			},
			API: api.URL,
		},
	), nil
}
