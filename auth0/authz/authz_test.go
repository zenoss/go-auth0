// +build integration

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
	}
	domain := getFromEnv("AUTH0_DOMAIN")
	s.authorization = auth0.AuthzClientFromCredentials(domain, api)
}

func TestAuthzTestSuite(t *testing.T) {
	suite.Run(t, new(AuthzTestSuite))
}
