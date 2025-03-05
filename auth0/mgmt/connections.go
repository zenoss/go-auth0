package mgmt

import (
	"github.com/google/go-querystring/query"
	"github.com/zenoss/go-auth0/auth0/http"
)

// ConnectionsService provides a service for connnection related functions
type ConnectionsService struct {
	c *http.Client
}

// Connection is a connection in Auth0
type Connection struct {
	Name               string         `json:"name,omitempty"`
	DisplayName        string         `json:"display_name,omitempty"`
	Options            map[string]any `json:"options,omitempty"`
	ID                 string         `json:"id,omitempty"`
	Strategy           string         `json:"strategy,omitempty"`
	Realms             []string       `json:"realms,omitempty"`
	EnabledClients     []string       `json:"enabled_clients,omitempty"`
	IsDomainConnection bool           `json:"is_domain_connection,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}

type ConnectionOpts struct {
	Name               string         `json:"name,omitempty"`
	DisplayName        string         `json:"display_name,omitempty"`
	Strategy           string         `json:"strategy,omitempty"`
	Options            map[string]any `json:"options,omitempty"`
	EnabledClients     []string       `json:"enabled_clients,omitempty"`
	IsDomainConnection bool           `json:"is_domain_connection,omitempty"`
	Realms             []string       `json:"realms,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}

type ConnectionUpdateOpts struct {
	DisplayName        string         `json:"display_name,omitempty"`
	Options            map[string]any `json:"options,omitempty"`
	EnabledClients     []string       `json:"enabled_clients,omitempty"`
	IsDomainConnection bool           `json:"is_domain_connection,omitempty"`
	Realms             []string       `json:"realms,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}
type ConnectionsPage struct {
	Start       int          `json:"start,omitempty"`
	Limit       int          `json:"limit,omitempty"`
	Length      int          `json:"length,omitempty"`
	Total       int          `json:"total,omitempty"`
	Connections []Connection `json:"connections,omitempty"`
}

// SearchConnectionsOpts defines what can be used to search connections
type SearchConnectionsOpts struct {
	PerPage       int      `url:"per_page,omitempty"`
	Page          int      `url:"page,omitempty"`
	IncludeTotals bool     `url:"include_totals,omitempty"`
	Strategy      []string `url:"strategy,omitempty"`
	Name          string   `url:"name,omitempty"`
	Fields        string   `url:"fields,omitempty"`
	IncludeFields bool     `url:"include_fields,omitempty"`
}

// Encode creates a url.Values encoding of SearchConnectionsOpts.
func (opts *SearchConnectionsOpts) Encode() (string, error) {
	vals, err := query.Values(opts)
	if err != nil {
		return "", err
	}
	return vals.Encode(), nil
}

// GetAll returns all connections
func (svc *ConnectionsService) GetAll() ([]Connection, error) {
	var connections []Connection
	err := svc.c.GetV2("/connections", &connections)
	return connections, err
}

// Get returns a connection
func (svc *ConnectionsService) Get(connectionID string) (Connection, error) {
	var connection Connection
	err := svc.c.Get("/connections/"+connectionID, &connection)
	return connection, err
}

// Search retrieves connections according to search criteria
func (svc *ConnectionsService) Search(opts SearchConnectionsOpts) (*ConnectionsPage, error) {
	var connectionsPage ConnectionsPage
	queryString, err := opts.Encode()
	if err != nil {
		return nil, err
	}
	url := "/connections"
	if queryString != "" {
		url = "/connections?" + queryString
	}
	if opts.IncludeTotals {
		err = svc.c.Get(url, &connectionsPage)
	} else {
		err = svc.c.Get(url, &connectionsPage.Connections)
	}
	return &connectionsPage, err
}

// Create creates a connection
func (svc *ConnectionsService) Create(opts ConnectionOpts) (Connection, error) {
	var connection Connection
	err := svc.c.Post("/connections", opts, &connection)
	return connection, err
}

// Delete deletes a connection
func (svc *ConnectionsService) Delete(connectionID string) error {
	return svc.c.Delete("/connections/"+connectionID, nil, nil)
}

func (svc *ConnectionsService) DeleteWithBody(connectionID string, body any) error {
	return svc.c.Delete("/connections/"+connectionID, &body, nil)
}

// Update updates a connection
func (svc *ConnectionsService) Update(connectionID string, opts ConnectionUpdateOpts) (Connection, error) {
	var connection Connection
	err := svc.c.Patch("/connections/"+connectionID, &opts, &connection)
	return connection, err
}
