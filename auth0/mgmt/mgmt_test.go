// +build integration

package mgmt_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/casbin/go-auth0/auth0"
)

type ManagementTestSuite struct {
	suite.Suite
	Client *auth0.Auth0
}

func (s *ManagementTestSuite) SetupSuite() {
	cfg := auth0.Config{
		ClientID:         os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret:     os.Getenv("AUTH0_CLIENT_SECRET"),
		Tenant:           os.Getenv("AUTH0_TENANT"),
		AuthorizationURL: os.Getenv("AUTH0_AUTHORIZATION_URL"),
	}
	client, err := cfg.ClientFromCredentials(os.Getenv("AUTH0_MANAGEMENT_API"))
	assert.Nil(s.T(), err)
	s.Client = client
}

func (suite *ManagementTestSuite) SetupTest() {
	cleanUpUsers(suite)
}

func TestManagementTestSuite(t *testing.T) {
	suite.Run(t, new(ManagementTestSuite))
}
