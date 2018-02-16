package auth0

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/zenoss/go-auth0/auth0/http"
)

// TokenService provides a service for token related functions
type TokenService struct {
	*http.Client
	API string
}

// TokenResponseBody contains token related information returned
// from a token request to Auth0
type TokenResponseBody struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    uint32 `json:"expires_in,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// TokenRequestBody contains fields that may be used as data in a POST
// to /oauth/token to request a token
type TokenRequestBody struct {
	GrantType    string `json:"grant_type,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Audience     string `json:"audience,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Scope        string `json:"scope,omitempty"`
	Code         string `json:"code,omitempty"`
	CodeVerifier string `json:"code_verifier,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	Realm        string `json:"realm,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// GetToken performs a generic call to /oauth/token using the body defined
// in a TokenRequestBody to get a TokenResponseBody, containing, at minimum,
// an access token, token type, and expiration
func (svc *TokenService) GetToken(body TokenRequestBody) (*TokenResponseBody, error) {
	var resBody TokenResponseBody
	reqBody, err := json.Marshal(body)
	if err != nil {
		return &resBody, errors.Wrap(err, "Cannot marshal TokenRequestBody")
	}
	err = svc.Post("/oauth/token", reqBody, &resBody)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot complete token request")
	}
	return &resBody, nil
}

// GetTokenFromClientCreds gets an access token to the target API
// using client credientials to authenticate
func (svc *TokenService) GetTokenFromClientCreds(clientID, clientSecret, API string) (*TokenResponseBody, error) {
	body := TokenRequestBody{
		GrantType:    "client_credentials",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     API,
	}
	return svc.GetToken(body)
}

// GetTokenFromUserPass gets an access token using
// username and password to authenticate
func (svc *TokenService) GetTokenFromUserPass(username, password, clientID string) (*TokenResponseBody, error) {
	body := TokenRequestBody{
		GrantType: "password",
		ClientID:  clientID,
		Username:  username,
		Password:  password,
	}
	return svc.GetToken(body)
}

// GetRealmTokenFromUserPass gets an access token using
// username and password to authenticate in a realm
func (svc *TokenService) GetRealmTokenFromUserPass(username, password, clientID, realm string) (*TokenResponseBody, error) {
	body := TokenRequestBody{
		GrantType: "http://auth0.com/oauth/grant-type/password-realm",
		ClientID:  clientID,
		Username:  username,
		Password:  password,
		Realm:     realm,
	}
	return svc.GetToken(body)
}
