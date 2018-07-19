// +build integration

package authz_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zenoss/go-auth0/auth0/authz"
)

const (
	groupName = "go-auth0-test-group"
	groupDesc = "A test group for go-auth0"
)

func createGroup(suite *AuthzTestSuite) authz.GroupStub {
	stub, err := suite.authorization.Groups.Create(groupName, groupDesc)
	assert.Nil(suite.T(), err)
	return stub
}

func getAllGroups(suite *AuthzTestSuite) []authz.Group {
	groups, err := suite.authorization.Groups.GetAll()
	assert.Nil(suite.T(), err)
	return groups
}

func deleteGroup(suite *AuthzTestSuite, ID string, ignoreErr bool) {
	err := suite.authorization.Groups.Delete(ID)
	if !ignoreErr {
		assert.Nil(suite.T(), err)
	}
}

func cleanUp(suite *AuthzTestSuite) {
	groups := getAllGroups(suite)
	var ID string
	for _, g := range groups {
		if g.Name == groupName {
			ID = g.ID
			break
		}
	}
	deleteGroup(suite, ID, true)
}

func (suite *AuthzTestSuite) SetupTest() {
	cleanUp(suite)
}

func (suite *AuthzTestSuite) TestGroupsCreateGetAllDelete() {
	t := suite.T()
	svc := suite.authorization.Groups
	// Create a group
	stub := createGroup(suite)
	assert.Equal(suite.T(), groupName, stub.Name)
	// Check we made it successfully
	group, err := svc.Get(stub.ID, true)
	assert.Nil(t, err)
	assert.Equal(t, groupName, group.Name)

	// And that getall shows it
	groups := getAllGroups(suite)
	found := false
	for _, g := range groups {
		if stub.ID == g.ID {
			found = true
		}
	}
	assert.True(t, found)

	// Update it
	group.Description = "go-auth0 test group"
	stub, err = svc.Update(group.GroupStub)
	assert.Nil(t, err)

	// Delete it
	deleteGroup(suite, stub.ID, false)
	// Check it was deleted
	groups = getAllGroups(suite)
	found = false
	for _, g := range groups {
		if stub.ID == g.ID {
			found = true
		}
	}
	assert.False(t, found)
}
