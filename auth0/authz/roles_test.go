// +build integration

package authz_test

import (
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/zenoss/go-auth0/auth0/authz"
)

const (
	roleName    = "go-auth0-test-role"
	roleDesc    = "A test role for go-auth0"
	roleAppType = "client"
)

var (
	roleAppID = os.Getenv("AUTH0_CLIENT_ID")
)

var rolePerms = []authz.Permission{}

func createRole(suite *AuthzTestSuite) authz.Role {
	role, err := suite.Client.Authz.Roles.Create(authz.Role{
		Name:            roleName,
		Description:     roleDesc,
		ApplicationType: roleAppType,
		ApplicationID:   roleAppID,
	})
	assert.Nil(suite.T(), err)
	return role
}

func getAllRoles(suite *AuthzTestSuite) []authz.Role {
	roles, err := suite.Client.Authz.Roles.GetAll()
	assert.Nil(suite.T(), err)
	return roles
}

func updateRole(suite *AuthzTestSuite, role authz.Role) authz.Role {
	role, err := suite.Client.Authz.Roles.Update(role)
	assert.Nil(suite.T(), err)
	return role
}

func deleteRole(suite *AuthzTestSuite, ID string, ignoreErr bool) {
	err := suite.Client.Authz.Roles.Delete(ID)
	if !ignoreErr {
		assert.Nil(suite.T(), err)
	}
}

func cleanUpRoles(suite *AuthzTestSuite) {
	roles := getAllRoles(suite)
	var ID string
	for _, p := range roles {
		if p.Name == roleName {
			ID = p.ID
			break
		}
	}
	if ID != "" {
		deleteRole(suite, ID, true)
	}
}

func (suite *AuthzTestSuite) TestRolesCreateGetAllDelete() {
	t := suite.T()
	svc := suite.Client.Authz.Roles
	// Create a role
	role := createRole(suite)
	assert.Equal(suite.T(), roleName, role.Name)
	// Check we made it successfully
	role, err := svc.Get(role.ID)
	assert.Nil(t, err)
	assert.Equal(t, roleName, role.Name)

	// The auth0 api won't let us update without including application type
	role.ApplicationType = roleAppType
	role.ApplicationID = roleAppID
	role.Description = "go-auth0 test role"
	role = updateRole(suite, role)

	// Delete it
	deleteRole(suite, role.ID, false)
	// Check it was deleted
	roles := getAllRoles(suite)
	found := false
	for _, p := range roles {
		if role.ID == p.ID {
			found = true
		}
	}
	assert.False(t, found)
}
