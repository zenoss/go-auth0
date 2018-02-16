package authz

import (
	"github.com/zenoss/go-auth0/auth0/http"
)

// PermissionsService provides a service for permission related functions
type PermissionsService struct {
	c *http.Client
}

// Permission is a permission
type Permission struct {
	ID              string `json:"_id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	ApplicationType string `json:"applicationType,omitempty"`
	ApplicationID   string `json:"applicationId,omitempty"`
}

// GetAll returns all permissions
func (svc *PermissionsService) GetAll() ([]Permission, error) {
	var permissions []Permission
	err := svc.c.Get("/api/permissions", &struct {
		Permissions *[]Permission `json:"permissions,omitempty"`
	}{Permissions: &permissions})
	return permissions, err
}

// Get returns a permissions
func (svc *PermissionsService) Get(ID string) (Permission, error) {
	var perm Permission
	err := svc.c.Get("/api/permissions/"+ID, &perm)
	return perm, err
}

// Create creates a permission
func (svc *PermissionsService) Create(perm Permission) (Permission, error) {
	var permResp Permission
	perm.ID = ""
	err := svc.c.Post("/api/permissions", &perm, &permResp)
	return permResp, err
}

// Delete deletes a permissions
func (svc *PermissionsService) Delete(ID string) error {
	return svc.c.Delete("/api/permissions/"+ID, nil, nil)
}

// Update creates a permission
func (svc *PermissionsService) Update(perm Permission) (Permission, error) {
	var permResp Permission
	permID := perm.ID
	perm.ID = ""
	err := svc.c.Put("/api/permissions/"+permID, &perm, &permResp)
	return permResp, err
}
