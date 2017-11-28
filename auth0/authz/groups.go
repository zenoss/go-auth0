package authz

// GroupsService provides a service for group related functions
type GroupsService struct {
	c Client
}

type groupStub struct {
	ID          string `json:"_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (stub *groupStub) ToGroupStub() *GroupStub {
	return &GroupStub{
		id:          stub.ID,
		Name:        stub.Name,
		Description: stub.Description,
	}
}

// GroupStub is a stub of a group
type GroupStub struct {
	id          string
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ID returns the id of the group
func (stub *GroupStub) ID() string {
	return stub.id
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
		item += "?expand"
	}
	err := svc.c.Get("/api/groups"+item, &group)
	return &group, err
}

// Create creates a group
func (svc *GroupsService) Create(stub *GroupStub) (*GroupStub, error) {
	var group groupStub
	err := svc.c.Post("/api/groups", Group{
		GroupStub: *stub,
	}, &group)
	return group.ToGroupStub(), err
}

// Delete deletes a groups
func (svc *GroupsService) Delete(ID string) error {
	return svc.c.Delete("/api/groups/"+ID, nil, nil)
}

// Update creates a group
func (svc *GroupsService) Update(stub *GroupStub) (*GroupStub, error) {
	var group groupStub
	err := svc.c.Put("/api/groups/"+stub.ID(), stub, &group)
	return group.ToGroupStub(), err
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
