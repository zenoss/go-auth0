package authz

import (
	"github.com/zenoss/go-auth0/auth0/http"
)

// RolesService provides a service for role related functions
type RolesService struct {
	c *http.Client
}

// Role is a role
type Role struct {
	ID              string   `json:"_id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	ApplicationType string   `json:"applicationType,omitempty"`
	ApplicationID   string   `json:"applicationId,omitempty"`
	PermissionIDs   []string `json:"permissions,omitempty"`
}

// GetAll returns all roles
func (svc *RolesService) GetAll() ([]Role, error) {
	var roles []Role
	err := svc.c.GetV2("/roles", &struct {
		Roles *[]Role `json:"roles,omitempty"`
	}{Roles: &roles})
	return roles, err
}

// Get returns a roles
func (svc *RolesService) Get(id string) (Role, error) {
	var r Role
	err := svc.c.Get("/roles/"+id, &r)
	return r, err
}

// Create creates a role
func (svc *RolesService) Create(r Role) (Role, error) {
	var roleResp Role
	r.ID = ""
	err := svc.c.Post("/roles", &r, &roleResp)
	return roleResp, err
}

// Delete deletes a roles
func (svc *RolesService) Delete(id string) error {
	return svc.c.Delete("/roles/"+id, nil, nil)
}

// Update creates a role
func (svc *RolesService) Update(r Role) (Role, error) {
	var roleResp Role
	roleID := r.ID
	r.ID = ""
	err := svc.c.Put("/roles/"+roleID, &r, &roleResp)
	return roleResp, err
}
