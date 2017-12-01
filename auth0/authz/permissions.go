package authz

// PermissionsService provides a service for permission related functions
type PermissionsService struct {
	c Client
}

type permission struct {
	ID              string `json:"_id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	ApplicationType string `json:"applicationType,omitempty"`
	ApplicationID   string `json:"applicationId,omitempty"`
}

func (stub *permission) ToPermission() *Permission {
	return &Permission{
		id:              stub.ID,
		Name:            stub.Name,
		Description:     stub.Description,
		ApplicationType: stub.ApplicationType,
		ApplicationID:   stub.ApplicationID,
	}
}

// Permission is a permission
type Permission struct {
	id              string
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	ApplicationType string `json:"applicationType,omitempty"`
	ApplicationID   string `json:"applicationId,omitempty"`
}

// ID returns the id of the permission
func (perm *Permission) ID() string {
	return perm.id
}

// GetAll returns all permissions
func (svc *PermissionsService) GetAll() (*[]Permission, error) {
	var permissions []permission
	err := svc.c.Get("/api/permissions", &struct {
		Permissions *[]permission `json:"permissions,omitempty"`
	}{Permissions: &permissions})
	perms := make([]Permission, len(permissions))
	for n, perm := range permissions {
		perms[n] = *perm.ToPermission()
	}
	return &perms, err
}

// Get returns a permissions
func (svc *PermissionsService) Get(ID string) (*Permission, error) {
	var perm permission
	err := svc.c.Get("/api/permissions/"+ID, &perm)
	return perm.ToPermission(), err
}

// Create creates a permission
func (svc *PermissionsService) Create(perm *Permission) (*Permission, error) {
	var permResp permission
	err := svc.c.Post("/api/permissions", perm, &permResp)
	return permResp.ToPermission(), err
}

// Delete deletes a permissions
func (svc *PermissionsService) Delete(ID string) error {
	return svc.c.Delete("/api/permissions/"+ID, nil, nil)
}

// Update creates a permission
func (svc *PermissionsService) Update(perm *Permission) (*Permission, error) {
	var permResp permission
	err := svc.c.Put("/api/permissions/"+perm.ID(), perm, &permResp)
	return permResp.ToPermission(), err
}
