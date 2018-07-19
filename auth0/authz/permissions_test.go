// +build integration

package authz_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zenoss/go-auth0/auth0/authz"
)

const (
	permName    = "go-auth0-test-permission"
	permDesc    = "A test permission for go-auth0"
	permAppType = "client"
)

var (
	permAppID = getFromEnv("AUTH0_AUTHORIZATION_CLIENT_ID")
)

func createPerm(suite *AuthzTestSuite) authz.Permission {
	perm, err := suite.authorization.Permissions.Create(authz.Permission{
		Name:            permName,
		Description:     permDesc,
		ApplicationType: permAppType,
		ApplicationID:   permAppID,
	})
	assert.Nil(suite.T(), err)
	return perm
}

func getAllPerms(suite *AuthzTestSuite) []authz.Permission {
	perms, err := suite.authorization.Permissions.GetAll()
	assert.Nil(suite.T(), err)
	return perms
}

func updatePerm(suite *AuthzTestSuite, perm authz.Permission) authz.Permission {
	perm, err := suite.authorization.Permissions.Update(perm)
	assert.Nil(suite.T(), err)
	return perm
}

func deletePerm(suite *AuthzTestSuite, ID string, ignoreErr bool) {
	err := suite.authorization.Permissions.Delete(ID)
	if !ignoreErr {
		assert.Nil(suite.T(), err)
	}
}

func cleanUpPerms(suite *AuthzTestSuite) {
	perms := getAllPerms(suite)
	var ID string
	for _, p := range perms {
		if p.Name == permName {
			ID = p.ID
			break
		}
	}
	if ID != "" {
		deletePerm(suite, ID, true)
	}
}

func (suite *AuthzTestSuite) TestPermsCreateGetAllDelete() {
	t := suite.T()
	svc := suite.authorization.Permissions
	// Create a permission
	permission := createPerm(suite)
	assert.Equal(suite.T(), permName, permission.Name)
	// Check we made it successfully
	permission, err := svc.Get(permission.ID)
	assert.Nil(t, err)
	assert.Equal(t, permName, permission.Name)

	// The auth0 api won't let us update without including application type
	permission.ApplicationType = permAppType
	permission.ApplicationID = permAppID
	permission.Description = "go-auth0 test permission"
	permission = updatePerm(suite, permission)

	// Delete it
	deletePerm(suite, permission.ID, false)
	// Check it was deleted
	permissions := getAllPerms(suite)
	found := false
	for _, p := range permissions {
		if permission.ID == p.ID {
			found = true
		}
	}
	assert.False(t, found)
}
