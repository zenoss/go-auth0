// +build integration

package auth0_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zenoss/go-auth0/auth0"
)

type Auth0TestSuite struct {
	suite.Suite
	Domain           string
	API              auth0.API
	ManagementAPI    auth0.API
	AuthorizationAPI auth0.API
}

func getFromEnv(envVar string) string {
	val := os.Getenv(envVar)
	if val == "" {
		panic("environment variable '" + envVar + "' must be set")
	}
	return val
}

func (s *Auth0TestSuite) SetupSuite() {
	s.Domain = getFromEnv("AUTH0_DOMAIN")
	s.API = auth0.API{
		URL:          getFromEnv("AUTH0_MANAGEMENT_API_URL"),
		Audience:     []string{getFromEnv("AUTH0_MANAGEMENT_API_AUDIENCE")},
		ClientID:     getFromEnv("AUTH0_MANAGEMENT_CLIENT_ID"),
		ClientSecret: getFromEnv("AUTH0_MANAGEMENT_CLIENT_SECRET"),
	}
	s.ManagementAPI = auth0.API{
		URL:          getFromEnv("AUTH0_MANAGEMENT_API_URL"),
		Audience:     []string{getFromEnv("AUTH0_MANAGEMENT_API_AUDIENCE")},
		ClientID:     getFromEnv("AUTH0_MANAGEMENT_CLIENT_ID"),
		ClientSecret: getFromEnv("AUTH0_MANAGEMENT_CLIENT_SECRET"),
	}
	s.AuthorizationAPI = auth0.API{
		URL:          getFromEnv("AUTH0_AUTHORIZATION_API_URL"),
		Audience:     []string{getFromEnv("AUTH0_AUTHORIZATION_API_AUDIENCE")},
		ClientID:     getFromEnv("AUTH0_AUTHORIZATION_CLIENT_ID"),
		ClientSecret: getFromEnv("AUTH0_AUTHORIZATION_CLIENT_SECRET"),
	}
}

func TestAuth0TestSuite(t *testing.T) {
	suite.Run(t, new(Auth0TestSuite))
}

func (suite *Auth0TestSuite) TestClientFromCredentials() {
	t := suite.T()
	// Create a go http client ready to make calls to an api
	c := auth0.ClientFromCredentials(suite.Domain, suite.API)
	// For this test, we make a call to the Auth0 management api,
	//   but this could be yours or anyone elses api
	ref := getFromEnv("AUTH0_MANAGEMENT_API_URL") + "clients"
	resp, err := c.Get(ref)
	assert.Nil(t, err)
	defer resp.Body.Close()
	clients, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	t.Logf("clients:\n%s\n", clients)
}

func (suite *Auth0TestSuite) TestMgmtClientFromCredentials() {
	t := suite.T()
	// Create a go http client ready to make calls to the auth0 mgmt api
	c := auth0.MgmtClientFromCredentials(suite.Domain, suite.ManagementAPI)
	users, err := c.Users.GetAll()
	assert.Nil(t, err)
	t.Logf("users:\n%v\n", users)
}

func (suite *Auth0TestSuite) TestAuthzClientFromCredentials() {
	t := suite.T()
	// Create a go http client ready to make calls to the auth0 authz api
	c := auth0.AuthzClientFromCredentials(suite.Domain, suite.AuthorizationAPI)
	groups, err := c.Groups.GetAll()
	assert.Nil(t, err)
	t.Logf("groups:\n%v\n", groups)
}
