package auth0

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

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
func (auth *Auth0) GetToken(body TokenRequestBody) (TokenResponseBody, error) {
	var resBody TokenResponseBody
	reqBody, err := json.Marshal(body)
	if err != nil {
		return resBody, errors.Wrap(err, "Cannot marshal TokenRequestBody")
	}
	req, err := http.NewRequest("POST", auth.site+"/oauth/token", bytes.NewBuffer(reqBody))
	if err != nil {
		return resBody, errors.Wrap(err, "Cannot create request")
	}
	resp, err := auth.do(req)
	if err != nil {
		return resBody, errors.Wrap(err, "Cannot complete token request")
	}
	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resBody, errors.Wrap(err, "Cannot read response body")
	}
	err = json.Unmarshal(bodyData, &resBody)
	if err != nil {
		return resBody, errors.Wrap(err, "Cannot unmarshal response into TokenResponseBody")
	}
	return resBody, nil
}

// GetTokenFromClientCreds gets an access token to the target API
// using client credientials to authenticate
func (auth *Auth0) GetTokenFromClientCreds(clientID, clientSecret, API string) (TokenResponseBody, error) {
	body := TokenRequestBody{
		GrantType:    "client_credentials",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     API,
	}
	return auth.GetToken(body)
}

// GetTokenFromUserPass gets an access token using
// username and password to authenticate
func (auth *Auth0) GetTokenFromUserPass(username, password, clientID string) (TokenResponseBody, error) {
	body := TokenRequestBody{
		GrantType: "password",
		ClientID:  clientID,
		Username:  username,
		Password:  password,
	}
	return auth.GetToken(body)
}

// GetRealmTokenFromUserPass gets an access token using
// username and password to authenticate in a realm
func (auth *Auth0) GetRealmTokenFromUserPass(username, password, clientID, realm string) (TokenResponseBody, error) {
	body := TokenRequestBody{
		GrantType: "http://auth0.com/oauth/grant-type/password-realm",
		ClientID:  clientID,
		Username:  username,
		Password:  password,
		Realm:     realm,
	}
	return auth.GetToken(body)
}
