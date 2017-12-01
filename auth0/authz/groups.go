package authz

import "github.com/pkg/errors"

// GroupsService provides a service for group related functions
type GroupsService struct {
	c Client
}

// GroupStub is a stub of a group
type GroupStub struct {
	ID          string `json:"_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Group is a group from the authorization extension
type Group struct {
	GroupStub
	Members  []string      `json:"members,omitempty"`
	Mappings []interface{} `json:"mappings,omitempty"`
}

// Mapping is a connection->group mapping for a group
type Mapping struct {
	GroupName      string `json:"groupName,omitempty"`
	ConnectionName string `json:"connectionName,omitempty"`
}

// GetAll returns all groups
func (svc *GroupsService) GetAll() (*[]Group, error) {
	var groups []Group
	err := svc.c.Get("/api/groups", &struct {
		Groups *[]Group `json:"groups,omitempty"`
	}{Groups: &groups})
	return &groups, err
}

// Get returns a groups
func (svc *GroupsService) Get(ID string, expand bool) (*Group, error) {
	var group Group
	item := "/" + ID
	if expand {
		item += "?expand=True"
	}
	err := svc.c.Get("/api/groups"+item, &group)
	return &group, err
}

// Create creates a group
func (svc *GroupsService) Create(name, description string) (*GroupStub, error) {
	var group GroupStub
	err := svc.c.Post("/api/groups", GroupStub{
		Name:        name,
		Description: description,
	}, &group)
	return &group, err
}

// Delete deletes a groups
func (svc *GroupsService) Delete(ID string) error {
	return svc.c.Delete("/api/groups/"+ID, nil, nil)
}

// Update updates a group
func (svc *GroupsService) Update(stub *GroupStub) (*GroupStub, error) {
	var group GroupStub
	err := svc.c.Put("/api/groups/"+stub.ID, stub, &group)
	return &group, err
}

// GetMappings get the mappings for a group
func (svc *GroupsService) GetMappings(ID string) (*[]Mapping, error) {
	var mappings []Mapping
	err := svc.c.Get("/api/groups/"+ID+"/mappings", &mappings)
	return &mappings, err
}

// CreateMappings creates one or more mappings for a group
func (svc *GroupsService) CreateMappings(ID string, mappings *[]Mapping) (*[]Mapping, error) {
	var mappingsResp []Mapping
	err := svc.c.Patch("/api/groups/"+ID+"/mappings", mappings, &mappingsResp)
	return &mappingsResp, err
}

// DeleteMappings creates one or more mappings for a group
func (svc *GroupsService) DeleteMappings(ID string, mappings *[]Mapping) (*[]Mapping, error) {
	var mappingsResp []Mapping
	err := svc.c.Delete("/api/groups/"+ID+"/mappings", mappings, &mappingsResp)
	return &mappingsResp, err
}

// GetMembers gets the members of a group
func (svc *GroupsService) GetMembers(ID string) (*[]string, error) {
	var members []string
	err := svc.c.Get("/api/groups/"+ID+"/members", &members)
	return &members, err
}

// AddMembers adds one or more members to a group
func (svc *GroupsService) AddMembers(ID string, members *[]string) (*[]string, error) {
	var membersResp []string
	err := svc.c.Patch("/api/groups/"+ID+"/members", members, &membersResp)
	return &membersResp, err
}

// DeleteMembers deletes one or more members from a group
func (svc *GroupsService) DeleteMembers(ID string, members *[]string) error {
	err := svc.c.Delete("/api/groups/"+ID+"/members", members, nil)
	return errors.Wrap(err, "go-auth0: cannot delete member from group")
}

// GetNestedMembers gets members in nested groups
func (svc *GroupsService) GetNestedMembers(ID string) (*[]string, error) {
	var members []string
	err := svc.c.Get("/api/groups/"+ID+"/members/nested", &members)
	return &members, err
}

// GetNestedGroups gets nested groups of a group
func (svc *GroupsService) GetNestedGroups(ID string) (*[]string, error) {
	var groups []string
	err := svc.c.Get("/api/groups/"+ID+"/nested", &groups)
	return &groups, err
}

// AddNestedGroups adds one or more nested groups to a group
func (svc *GroupsService) AddNestedGroups(ID string, groups *[]string) (*[]string, error) {
	var groupsResp []string
	err := svc.c.Patch("/api/groups/"+ID+"/nested", groups, &groupsResp)
	return &groupsResp, err
}

// DeleteNestedGroups deletes one or more nested groups from a group
func (svc *GroupsService) DeleteNestedGroups(ID string, groups *[]string) (*[]string, error) {
	var groupsResp []string
	err := svc.c.Delete("/api/groups/"+ID+"/nested", groups, &groupsResp)
	return &groupsResp, err
}

// GetGroupRoles gets the roles for a groups
func (svc *GroupsService) GetGroupRoles(ID string) (*[]string, error) {
	var roles []string
	err := svc.c.Get("/api/groups/"+ID+"/roles", &roles)
	return &roles, err
}

// AddGroupRoles adds one or more roles to a group
func (svc *GroupsService) AddGroupRoles(ID string, roles *[]string) (*[]string, error) {
	var rolesResp []string
	err := svc.c.Patch("/api/groups/"+ID+"/roles", roles, &rolesResp)
	return &rolesResp, err
}

// DeleteGroupRoles deletes one or more roles from a group
func (svc *GroupsService) DeleteGroupRoles(ID string, roles *[]string) (*[]string, error) {
	var rolesResp []string
	err := svc.c.Delete("/api/groups/"+ID+"/roles", roles, &rolesResp)
	return &rolesResp, err
}

// GetNestedRoles gets roles of nested groups from a group
func (svc *GroupsService) GetNestedRoles(ID string) (*[]string, error) {
	var roles []string
	err := svc.c.Get("/api/groups/"+ID+"/roles/nested", &roles)
	return &roles, err
}
