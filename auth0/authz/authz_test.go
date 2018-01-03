// +build integration

package authz_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zenoss/go-auth0/auth0"
)

type AuthzTestSuite struct {
	suite.Suite
	Client *auth0.Auth0
}

func (s *AuthzTestSuite) SetupSuite() {
	cfg := auth0.Config{
		ClientID:         os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret:     os.Getenv("AUTH0_CLIENT_SECRET"),
		Tenant:           os.Getenv("AUTH0_TENANT"),
		AuthorizationURL: os.Getenv("AUTH0_AUTHORIZATION_URL"),
	}
	client, err := cfg.ClientFromCredentials(os.Getenv("AUTH0_AUTHORIZATION_API"))
	assert.Nil(s.T(), err)
	s.Client = client
}

func TestAuthzTestSuite(t *testing.T) {
	suite.Run(t, new(AuthzTestSuite))
}
