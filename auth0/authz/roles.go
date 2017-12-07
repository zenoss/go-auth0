package authz

// RolesService provides a service for role related functions
type RolesService struct {
	c Client
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
	err := svc.c.Get("/api/roles", &struct {
		Roles []Role `json:"roles,omitempty"`
	}{Roles: roles})
	return roles, err
}

// Get returns a roles
func (svc *RolesService) Get(ID string) (Role, error) {
	var r Role
	err := svc.c.Get("/api/roles/"+ID, &r)
	return r, err
}

// Create creates a role
func (svc *RolesService) Create(r Role) (Role, error) {
	var roleResp Role
	r.ID = ""
	err := svc.c.Post("/api/roles", &r, &roleResp)
	return roleResp, err
}

// Delete deletes a roles
func (svc *RolesService) Delete(ID string) error {
	return svc.c.Delete("/api/roles/"+ID, nil, nil)
}

// Update creates a role
func (svc *RolesService) Update(r Role) (Role, error) {
	var roleResp Role
	roleID := r.ID
	r.ID = ""
	err := svc.c.Put("/api/roles/"+roleID, &r, &roleResp)
	return roleResp, err
}
