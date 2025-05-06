package authz

import (
	"fmt"

	"github.com/zenoss/go-auth0/auth0/http"
)

// GroupsService provides a service for group related functions
type GroupsService struct {
	c *http.Client
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
	Members  []string  `json:"members,omitempty"`
	Mappings []Mapping `json:"mappings,omitempty"`
	Roles    []string  `json:"roles,omitempty"`
}

// Mapping is a connection->group mapping for a group
type Mapping struct {
	ID             string `json:"_id,omitempty"`
	GroupName      string `json:"groupName,omitempty"`
	ConnectionName string `json:"connectionName,omitempty"`
}

// GetAll returns all groups
func (svc *GroupsService) GetAll() ([]Group, error) {
	var groups []Group
	err := svc.c.Get("/groups", &struct {
		Groups *[]Group `json:"groups,omitempty"`
	}{Groups: &groups})
	return groups, err
}

// Get returns a groups
func (svc *GroupsService) Get(groupID string, expand bool) (Group, error) {
	var group Group
	item := "/" + groupID
	if expand {
		item += "?expand=True"
	}
	err := svc.c.Get("/groups"+item, &group)
	return group, err
}

// Create creates a group
func (svc *GroupsService) Create(name, description string) (GroupStub, error) {
	var group GroupStub
	err := svc.c.Post("/groups", GroupStub{
		Name:        name,
		Description: description,
	}, &group)
	return group, err
}

// Delete deletes a groups
func (svc *GroupsService) Delete(groupID string) error {
	return svc.c.Delete("/groups/"+groupID, nil, nil)
}

// Update updates a group
func (svc *GroupsService) Update(stub GroupStub) (GroupStub, error) {
	var group GroupStub
	stubID := stub.ID
	stub.ID = ""
	err := svc.c.Put("/groups/"+stubID, &stub, &group)
	return group, err
}

// GetMappings get the mappings for a group
func (svc *GroupsService) GetMappings(groupID string) ([]Mapping, error) {
	var mappings []Mapping
	err := svc.c.Get("/groups/"+groupID+"/mappings", &mappings)
	return mappings, err
}

// CreateMappings creates one or more mappings for a group
func (svc *GroupsService) CreateMappings(groupID string, mappings []Mapping) error {
	for _, mapping := range mappings {
		mapping.ID = ""
	}
	err := svc.c.Patch("/groups/"+groupID+"/mappings", mappings, nil)
	return err
}

// DeleteMappings creates one or more mappings for a group
func (svc *GroupsService) DeleteMappings(groupID string, mappingIDs []string) error {
	err := svc.c.Delete("/groups/"+groupID+"/mappings", mappingIDs, nil)
	return err
}

// GetMembers gets the members of a group
func (svc *GroupsService) GetMembers(groupID string) ([]string, error) {
	var members []string
	err := svc.c.Get("/groups/"+groupID+"/members", &members)
	return members, err
}

// AddMembers adds one or more members to a group
func (svc *GroupsService) AddMembers(groupID string, members []string) ([]string, error) {
	var membersResp []string
	err := svc.c.Patch("/groups/"+groupID+"/members", members, &membersResp)
	return membersResp, err
}

// DeleteMembers deletes one or more members from a group
func (svc *GroupsService) DeleteMembers(groupID string, members []string) error {
	err := svc.c.Delete("/groups/"+groupID+"/members", members, nil)
	if err != nil {
		return fmt.Errorf("go-auth0: cannot delete members from group: %w", err)
	}
	return err
}

// GetNestedMembers gets members in nested groups
func (svc *GroupsService) GetNestedMembers(groupID string) ([]string, error) {
	var members []string
	err := svc.c.Get("/groups/"+groupID+"/members/nested", &members)
	return members, err
}

// GetNestedGroups gets nested groups of a group
func (svc *GroupsService) GetNestedGroups(groupID string) ([]string, error) {
	var groups []string
	err := svc.c.Get("/groups/"+groupID+"/nested", &groups)
	return groups, err
}

// AddNestedGroups adds one or more nested groups to a group
func (svc *GroupsService) AddNestedGroups(groupID string, groups []string) ([]string, error) {
	var groupsResp []string
	err := svc.c.Patch("/groups/"+groupID+"/nested", groups, &groupsResp)
	return groupsResp, err
}

// DeleteNestedGroups deletes one or more nested groups from a group
func (svc *GroupsService) DeleteNestedGroups(groupID string, groups []string) ([]string, error) {
	var groupsResp []string
	err := svc.c.Delete("/groups/"+groupID+"/nested", groups, &groupsResp)
	return groupsResp, err
}

// GetGroupRoles gets the roles for a groups
func (svc *GroupsService) GetGroupRoles(groupID string) ([]string, error) {
	var roles []string
	err := svc.c.Get("/groups/"+groupID+"/roles", &roles)
	return roles, err
}

// AddGroupRoles adds one or more roles to a group
func (svc *GroupsService) AddGroupRoles(groupID string, roles []string) error {
	err := svc.c.Patch("/groups/"+groupID+"/roles", roles, nil)
	return err
}

// DeleteGroupRoles deletes one or more roles from a group
func (svc *GroupsService) DeleteGroupRoles(groupID string, roles []string) error {
	err := svc.c.Delete("/groups/"+groupID+"/roles", roles, nil)
	return err
}

// GetNestedRoles gets roles of nested groups from a group
func (svc *GroupsService) GetNestedRoles(groupID string) ([]string, error) {
	var roles []string
	err := svc.c.Get("/groups/"+groupID+"/roles/nested", &roles)
	return roles, err
}
