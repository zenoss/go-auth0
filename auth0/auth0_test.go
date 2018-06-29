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

func (s *Auth0TestSuite) SetupSuite() {
	s.Domain = os.Getenv("AUTH0_DOMAIN")
	s.API = auth0.API{
		URL:          os.Getenv("AUTH0_MANAGEMENT_API_URL"),
		Audience:     []string{os.Getenv("AUTH0_MANAGEMENT_API_AUDIENCE")},
		ClientID:     os.Getenv("AUTH0_MANAGEMENT_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_MANAGEMENT_CLIENT_SECRET"),
	}
	// auth0.API{
	// 	URL:          os.Getenv("AUTH0_MYAPI_API_URL"),
	// 	Audience:     []string{os.Getenv("AUTH0_MYAPI_API_AUDIENCE")},
	// 	ClientID:     os.Getenv("AUTH0_MYAPI_CLIENT_ID"),
	// 	ClientSecret: os.Getenv("AUTH0_MYAPI_CLIENT_SECRET"),
	// }
	s.ManagementAPI = auth0.API{
		URL:          os.Getenv("AUTH0_MANAGEMENT_API_URL"),
		Audience:     []string{os.Getenv("AUTH0_MANAGEMENT_API_AUDIENCE")},
		ClientID:     os.Getenv("AUTH0_MANAGEMENT_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_MANAGEMENT_CLIENT_SECRET"),
	}
	s.AuthorizationAPI = auth0.API{
		URL:          os.Getenv("AUTH0_AUTHORIZATION_API_URL"),
		Audience:     []string{os.Getenv("AUTH0_AUTHORIZATION_API_AUDIENCE")},
		ClientID:     os.Getenv("AUTH0_AUTHORIZATION_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_AUTHORIZATION_CLIENT_SECRET"),
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
	ref := os.Getenv("AUTH0_MANAGEMENT_API_URL") + "clients"
	t.Logf("ref: %s", ref)
	resp, err := c.Get(ref)
	assert.Nil(t, err)
	defer resp.Body.Close()
	clients, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	t.Logf("resp:\n%s\n", clients)
}

func (suite *Auth0TestSuite) TestMgmtClientFromCredentials() {
	t := suite.T()
	// Create a go http client ready to make calls to the auth0 mgmt api
	c := auth0.MgmtClientFromCredentials(suite.Domain, suite.API)
	users, err := c.Users.GetAll()
	assert.Nil(t, err)
	t.Logf("users:\n%v\n", users)
}

func (suite *Auth0TestSuite) TestAuthzClientFromCredentials() {
	t := suite.T()
	// Create a go http client ready to make calls to the auth0 mgmt api
	c := auth0.AuthzClientFromCredentials(suite.Domain, suite.API)
	groups, err := c.Groups.GetAll()
	assert.Nil(t, err)
	t.Logf("groups:\n%v\n", groups)
}
