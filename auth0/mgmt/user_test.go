// +build integration

package mgmt_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/casbin/go-auth0/auth0/mgmt"
)

const (
	userEmail = "someone@example.com"
	userPass  = "wildwest"
)

func createUser(suite *ManagementTestSuite) mgmt.User {
	user, err := suite.Client.Mgmt.Users.Create(mgmt.UserOpts{
		Connection: "Username-Password-Authentication",
		Email:      userEmail,
		Password:   userPass,
	})
	assert.Nil(suite.T(), err)
	return user
}

func getAllUsers(suite *ManagementTestSuite) []mgmt.User {
	users, err := suite.Client.Mgmt.Users.GetAll()
	assert.Nil(suite.T(), err)
	return users
}

func deleteUser(suite *ManagementTestSuite, ID string, ignoreErr bool) {
	err := suite.Client.Mgmt.Users.Delete(ID)
	if !ignoreErr {
		assert.Nil(suite.T(), err)
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
	svc := suite.Client.Mgmt.Users

	// Create a user
	user := createUser(suite)
	assert.Equal(t, userEmail, user.Email)

	// Check we made it successfully
	user, err := svc.Get(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, userEmail, user.Email)

	// Update it
	update := mgmt.UserUpdateOpts{
		Blocked: true,
	}
	user, err = svc.Update(user.ID, update)
	assert.Nil(t, err)
	assert.Equal(t, true, user.Blocked)

	// Delete it
	deleteUser(suite, user.ID, false)

	// Check it was deleted
	_, err = svc.Get(user.ID)
	assert.NotNil(t, err)
}
