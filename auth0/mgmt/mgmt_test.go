// +build integration

// These integration tests run against an actual instance of Auth0,
//  but should not be run against a production tenant.

// The tests in this package will require that the management api is enabled
//  and that the following environment variables are defined
//  - AUTH0_MANAGEMENT_API_URL
//  - AUTH0_MANAGEMENT_API_AUDIENCE
//  - AUTH0_MANAGEMENT_CLIENT_ID
//  - AUTH0_MANAGEMENT_CLIENT_SECRET
// API URL and API Audience are likely the same and should be under API at auth0 dashboard
// The Client can be any enabled for the mgmt api and Username-Password-Authentication connection

package mgmt_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zenoss/go-auth0/auth0"
	"github.com/zenoss/go-auth0/auth0/mgmt"
)

type ManagementTestSuite struct {
	suite.Suite
	management *mgmt.ManagementService
}

func getFromEnv(envVar string) string {
	val := os.Getenv(envVar)
	if val == "" {
		panic("environment variable '" + envVar + "' must be set")
	}
	return val
}

func (s *ManagementTestSuite) SetupSuite() {
	api := auth0.API{
		URL:          getFromEnv("AUTH0_MANAGEMENT_API_URL"),
		Audience:     []string{getFromEnv("AUTH0_MANAGEMENT_API_AUDIENCE")},
		ClientID:     getFromEnv("AUTH0_MANAGEMENT_CLIENT_ID"),
		ClientSecret: getFromEnv("AUTH0_MANAGEMENT_CLIENT_SECRET"),
	}
	domain := getFromEnv("AUTH0_DOMAIN")
	s.management = auth0.MgmtClientFromCredentials(domain, api)
}

func (suite *ManagementTestSuite) SetupTest() {
	cleanUpUsers(suite)
}

func TestManagementTestSuite(t *testing.T) {
	suite.Run(t, new(ManagementTestSuite))
}
