package authz

// RolesService provides a service for role related functions
type RolesService struct {
	c Client
}

type role struct {
	ID              string       `json:"_id,omitempty"`
	Name            string       `json:"name,omitempty"`
	Description     string       `json:"description,omitempty"`
	ApplicationType string       `json:"applicationType,omitempty"`
	ApplicationID   string       `json:"applicationId,omitempty"`
	Permissions     []Permission `json:"permissions,omitempty"`
}

func (stub *role) ToRole() *Role {
	return &Role{
		id:              stub.ID,
		Name:            stub.Name,
		Description:     stub.Description,
		ApplicationType: stub.ApplicationType,
		ApplicationID:   stub.ApplicationID,
		Permissions:     stub.Permissions,
	}
}

// Role is a role
type Role struct {
	id              string
	Name            string       `json:"name,omitempty"`
	Description     string       `json:"description,omitempty"`
	ApplicationType string       `json:"applicationType,omitempty"`
	ApplicationID   string       `json:"applicationId,omitempty"`
	Permissions     []Permission `json:"permissions,omitempty"`
}

// ID returns the id of the role
func (r *Role) ID() string {
	return r.id
}

// GetAll returns all roles
func (svc *RolesService) GetAll() (*[]Role, error) {
	var rolesResp []role
	err := svc.c.Get("/api/roles", &struct {
		Roles *[]role `json:"roles,omitempty"`
	}{Roles: &rolesResp})
	roles := make([]Role, len(rolesResp))
	for n, r := range rolesResp {
		roles[n] = *r.ToRole()
	}
	return &roles, err
}

// Get returns a roles
func (svc *RolesService) Get(ID string) (*Role, error) {
	var r role
	err := svc.c.Get("/api/roles/"+ID, &r)
	return r.ToRole(), err
}

// Create creates a role
func (svc *RolesService) Create(r *Role) (*Role, error) {
	var roleResp role
	err := svc.c.Post("/api/roles", r, &roleResp)
	return roleResp.ToRole(), err
}

// Delete deletes a roles
func (svc *RolesService) Delete(ID string) error {
	return svc.c.Delete("/api/roles/"+ID, nil, nil)
}

// Update creates a role
func (svc *RolesService) Update(r *Role) (*Role, error) {
	var roleResp role
	err := svc.c.Put("/api/roles/"+r.ID(), r, &roleResp)
	return roleResp.ToRole(), err
}
