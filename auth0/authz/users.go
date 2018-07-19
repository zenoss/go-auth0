package authz

import (
	"github.com/pkg/errors"
	"github.com/zenoss/go-auth0/auth0/http"
)

// UsersService provides a service for user related functions
type UsersService struct {
	c *http.Client
}

// User is a user
type User struct {
	ID     string  `json:"_id,omitempty"`
	Roles  []Role  `json:"roles,omitempty"`
	Groups []Group `json:"groups,omitempty"`
}

// GetGroups returns the groups for a user
func (svc *UsersService) GetGroups(ID string, expand bool) ([]GroupStub, error) {
	var groupResp []GroupStub
	item := "/" + ID + "/groups"
	if expand {
		item += "?expand"
	}
	err := svc.c.Get("/users"+item, &groupResp)
	if err != nil {
		return nil, errors.Wrap(err, "go-auth0: cannot get groups for user")
	}
	return groupResp, err
}

// AddGroups puts the user in one or more groups
func (svc *UsersService) AddGroups(ID string, groups []string) error {
	err := svc.c.Patch("/users/"+ID+"/groups", &groups, nil)
	if err != nil {
		return errors.Wrap(err, "go-auth0: cannot add groups for user")
	}
	return nil
}

// GetAllGroups returns the groups for a user including nested groups
func (svc *UsersService) GetAllGroups(ID string) ([]GroupStub, error) {
	var groupResp []GroupStub
	err := svc.c.Get("/users/"+ID+"/groups/calculate", &groupResp)
	if err != nil {
		return nil, errors.Wrap(err, "go-auth0: cannot get groups for user")
	}
	return groupResp, err
}

// GetRoles returns the roles for a user
func (svc *UsersService) GetRoles(ID string) ([]Role, error) {
	var roleResp []Role
	err := svc.c.Get("/users/"+ID+"/roles", &roleResp)
	if err != nil {
		return nil, errors.Wrap(err, "go-auth0: cannot get roles for user")
	}
	roles := make([]Role, len(roleResp))
	for n, r := range roleResp {
		roles[n] = r
	}
	return roles, err
}

// AddRoles gives the user one or more roles
func (svc *UsersService) AddRoles(ID string, roles []string) error {
	err := svc.c.Patch("/users/"+ID+"/roles", &roles, nil)
	if err != nil {
		return errors.Wrap(err, "go-auth0: cannot add roles for user")
	}
	return nil
}

// RemoveRoles removes one or more roles from the user
func (svc *UsersService) RemoveRoles(ID string, roles []string) error {
	err := svc.c.Delete("/users/"+ID+"/roles", &roles, nil)
	if err != nil {
		return errors.Wrap(err, "go-auth0: cannot add roles for user")
	}
	return nil
}

// GetAllRoles returns all roles for a user, including through group membership
func (svc *UsersService) GetAllRoles(ID string) ([]Role, error) {
	var roleResp []Role
	err := svc.c.Get("/users/"+ID+"/roles/calculate", &roleResp)
	if err != nil {
		return nil, errors.Wrap(err, "go-auth0: cannot get all roles for user")
	}
	roles := make([]Role, len(roleResp))
	for n, r := range roleResp {
		roles[n] = r
	}
	return roles, err
}

// ExecAuthPolicy executes the authorization policy for a user in the context of a client
func (svc *UsersService) ExecAuthPolicy(ID, policyID, connection string, groups []string) error {
	body := struct {
		ConnectionName string   `json:"connectionName,omitempty"`
		Groups         []string `json:"groups,omitempty"`
	}{
		ConnectionName: connection,
		Groups:         groups,
	}
	return svc.c.Post("/users/"+ID+"/policy/"+policyID, body, nil)
}
