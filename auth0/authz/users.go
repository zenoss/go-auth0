package authz

import (
	"fmt"

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
func (svc *UsersService) GetGroups(id string, expand bool) ([]GroupStub, error) {
	var groupResp []GroupStub
	item := "/" + id + "/groups"
	if expand {
		item += "?expand"
	}
	err := svc.c.Get("/users"+item, &groupResp)
	if err != nil {
		return nil, fmt.Errorf("go-auth0: cannot get groups for user: %w", err)
	}
	return groupResp, err
}

// AddGroups puts the user in one or more groups
func (svc *UsersService) AddGroups(id string, groups []string) error {
	err := svc.c.Patch("/users/"+id+"/groups", &groups, nil)
	if err != nil {
		return fmt.Errorf("go-auth0: cannot add groups for user: %w", err)
	}
	return nil
}

// GetAllGroups returns the groups for a user including nested groups
func (svc *UsersService) GetAllGroups(id string) ([]GroupStub, error) {
	var groupResp []GroupStub
	err := svc.c.Get("/users/"+id+"/groups/calculate", &groupResp)
	if err != nil {
		return nil, fmt.Errorf("go-auth0: cannot get all groups for user: %w", err)
	}
	return groupResp, err
}

// GetRoles returns the roles for a user
func (svc *UsersService) GetRoles(id string) ([]Role, error) {
	var roleResp []Role
	err := svc.c.Get("/users/"+id+"/roles", &roleResp)
	if err != nil {
		return nil, fmt.Errorf("go-auth0: cannot get roles for user: %w", err)
	}
	roles := make([]Role, len(roleResp))
	copy(roles, roleResp)
	return roles, err
}

// AddRoles gives the user one or more roles
func (svc *UsersService) AddRoles(id string, roles []string) error {
	err := svc.c.Patch("/users/"+id+"/roles", &roles, nil)
	if err != nil {
		return fmt.Errorf("go-auth0: cannot add roles for user: %w", err)
	}
	return nil
}

// RemoveRoles removes one or more roles from the user
func (svc *UsersService) RemoveRoles(id string, roles []string) error {
	err := svc.c.Delete("/users/"+id+"/roles", &roles, nil)
	if err != nil {
		return fmt.Errorf("go-auth0: cannot remove roles for user: %w", err)
	}
	return nil
}

// GetAllRoles returns all roles for a user, including through group membership
func (svc *UsersService) GetAllRoles(id string) ([]Role, error) {
	var roleResp []Role
	err := svc.c.Get("/users/"+id+"/roles/calculate", &roleResp)
	if err != nil {
		return nil, fmt.Errorf("go-auth0: cannot get all roles for user: %w", err)
	}
	roles := make([]Role, len(roleResp))
	copy(roles, roleResp)
	return roles, err
}

// ExecAuthPolicy executes the authorization policy for a user in the context of a client
func (svc *UsersService) ExecAuthPolicy(id, policyID, connection string, groups []string) error {
	body := struct {
		ConnectionName string   `json:"connectionName,omitempty"`
		Groups         []string `json:"groups,omitempty"`
	}{
		ConnectionName: connection,
		Groups:         groups,
	}
	return svc.c.Post("/users/"+id+"/policy/"+policyID, body, nil)
}
