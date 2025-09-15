//go:build integration
// +build integration

package mgmt_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/zenoss/go-auth0/auth0/mgmt"
)

const (
	userEmail = "someone@example.com"
	userPass  = "WildWest000"
)

func createUser(suite *ManagementTestSuite) mgmt.User {
	user, err := suite.management.Users.Create(mgmt.UserOpts{
		Connection: "Username-Password-Authentication",
		Email:      userEmail,
		Password:   userPass,
	})
	suite.NoError(err)

	return user
}

func getAllUsers(suite *ManagementTestSuite) []mgmt.User {
	users, err := suite.management.Users.GetAll()
	suite.NoError(err)

	return users
}

func searchUsers(suite *ManagementTestSuite, searchOpts mgmt.SearchUsersOpts) *mgmt.UsersPage {
	usersPage, err := suite.management.Users.Search(searchOpts)
	suite.NoError(err)

	return usersPage
}

func deleteUser(suite *ManagementTestSuite, ID string, ignoreErr bool) {
	err := suite.management.Users.Delete(ID)
	if !ignoreErr {
		suite.NoError(err)
	}
}

func cleanUpUsers(suite *ManagementTestSuite) {
	users := getAllUsers(suite)

	var ID string

	for _, user := range users {
		if user.Email == userEmail {
			ID = user.ID
			break
		}
	}

	if ID != "" {
		deleteUser(suite, ID, true)
	}
}

func (suite *ManagementTestSuite) TestUsersCreateGetAllDelete() {
	t := suite.T()
	svc := suite.management.Users

	// Check if user existed before the test and remove
	cleanUpUsers(suite)

	// Create a user
	user := createUser(suite)
	assert.Equal(t, userEmail, user.Email)

	// Check we made it successfully
	user, err := svc.Get(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, userEmail, user.Email)

	// Check that we can search
	searchOpts := mgmt.SearchUsersOpts{
		Q: fmt.Sprintf(`email:"%s"`, userEmail),
	}
	users, err := svc.Search(searchOpts)
	assert.NoError(t, err)
	assert.NotNil(t, users)

	// Update it
	update := mgmt.UserUpdateOpts{
		Blocked: true,
	}
	user, err = svc.Update(user.ID, update)
	assert.NoError(t, err)
	assert.True(t, user.Blocked)

	// Delete it
	deleteUser(suite, user.ID, false)

	// Check it was deleted
	_, err = svc.Get(user.ID)
	assert.Error(t, err)
}
