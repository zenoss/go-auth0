// +build integration

// These integration tests run against an actual instance of Auth0,
//  but should not be run against a production tenant.

// The tests in this package will require that the authorization extension is installed,
//  and that the following environment variables are defined
//  - AUTH0_AUTHORIZATION_API_URL
//  - AUTH0_AUTHORIZATION_API_AUDIENCE
//  - AUTH0_AUTHORIZATION_CLIENT_ID
//  - AUTH0_AUTHORIZATION_CLIENT_SECRET
// API URL and API Audience are found by opening auth0 authorization, click your profile -> API
// The Client can be any enabled for the authorization api

package authz_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zenoss/go-auth0/auth0"
	"github.com/zenoss/go-auth0/auth0/authz"
)

type AuthzTestSuite struct {
	suite.Suite
	authorization *authz.AuthorizationService
}

func getFromEnv(envVar string) string {
	val := os.Getenv(envVar)
	if val == "" {
		panic("environment variable '" + envVar + "' must be set")
	}
	return val
}

func (s *AuthzTestSuite) SetupSuite() {
	api := auth0.API{
		URL:          getFromEnv("AUTH0_AUTHORIZATION_API_URL"),
		Audience:     []string{getFromEnv("AUTH0_AUTHORIZATION_API_AUDIENCE")},
		ClientID:     getFromEnv("AUTH0_AUTHORIZATION_CLIENT_ID"),
		ClientSecret: getFromEnv("AUTH0_AUTHORIZATION_CLIENT_SECRET"),
		Scopes:       []string{"delete:permissions", "delete:roles"},
	}
	domain := getFromEnv("AUTH0_DOMAIN")
	s.authorization = auth0.AuthzClientFromCredentials(domain, api)
}

func TestAuthzTestSuite(t *testing.T) {
	suite.Run(t, new(AuthzTestSuite))
}
