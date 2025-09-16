//go:build integration
// +build integration

package authz_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zenoss/go-auth0/auth0/authz"
)

const (
	roleName    = "go-auth0-test-role"
	roleDesc    = "A test role for go-auth0"
	roleAppType = "client"
)

var roleAppID = getFromEnv("AUTH0_AUTHORIZATION_CLIENT_ID")

var rolePerms = []authz.Permission{}

func createRole(suite *AuthzTestSuite) authz.Role {
	role, err := suite.authorization.Roles.Create(authz.Role{
		Name:            roleName,
		Description:     roleDesc,
		ApplicationType: roleAppType,
		ApplicationID:   roleAppID,
	})
	suite.NoError(err)

	return role
}

func getAllRoles(suite *AuthzTestSuite) []authz.Role {
	roles, err := suite.authorization.Roles.GetAll()
	suite.NoError(err)

	return roles
}

func updateRole(suite *AuthzTestSuite, role authz.Role) authz.Role {
	role, err := suite.authorization.Roles.Update(role)
	suite.NoError(err)

	return role
}

func deleteRole(suite *AuthzTestSuite, ID string, ignoreErr bool) {
	err := suite.authorization.Roles.Delete(ID)
	if !ignoreErr {
		suite.NoError(err)
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

	// Check if role existed before the test and remove
	cleanUpRoles(suite)

	svc := suite.authorization.Roles
	// Create a role
	role := createRole(suite)
	suite.Equal(roleName, role.Name)
	// Check we made it successfully
	role, err := svc.Get(role.ID)
	assert.NoError(t, err)
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
