package mgmt

import (
    "fmt"
    "github.com/google/go-querystring/query"
    "github.com/zenoss/go-auth0/auth0/http"
)

// UsersService provides a service for user related functions
type UsersService struct {
	c *http.Client
}

// User is a user in Auth0
type User struct {
	Email         string                 `json:"email,omitempty"`
	EmailVerified bool                   `json:"email_verified,omitempty"`
	Username      string                 `json:"username,omitempty"`
	PhoneNumber   string                 `json:"phone_number,omitempty"`
	PhoneVerified bool                   `json:"phone_verified,omitempty"`
	ID            string                 `json:"user_id,omitempty"`
	CreatedAt     string                 `json:"created_at,omitempty"`
	UpdatedAt     string                 `json:"updated_at,omitempty"`
	Identities    []Identity             `json:"identities,omitempty"`
	AppMetadata   map[string]interface{} `json:"app_metadata,omitempty"`
	UserMetadata  map[string]interface{} `json:"user_metadata,omitempty"`
	Picture       string                 `json:"picture,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Nickname      string                 `json:"nickname,omitempty"`
	Multifactor   []string               `json:"multifactor,omitempty"`
	LastIP        string                 `json:"last_ip,omitempty"`
	LastLogin     string                 `json:"last_login,omitempty"`
	LoginCount    uint                   `json:"logins_count,omitempty"`
	Blocked       bool                   `json:"blocked,omitempty"`
	FirstName     string                 `json:"given_name,omitempty"`
	LastName      string                 `json:"family_name,omitempty"`
}

type UsersPage struct {
    Start  int `json:"start,omitempty"`
    Limit  int `json:"limit,omitempty"`
    Length int `json:"length,omitempty"`
    Total int `json:"total,omitempty"`
    Users []User `json:"users,omitempty"`
}

// UserOpts are options which can be used to create a User
type UserOpts struct {
	ID            string                 `json:"user_id,omitempty"`
	Connection    string                 `json:"connection,omitempty"`
	Email         string                 `json:"email,omitempty"`
	Username      string                 `json:"username,omitempty"`
	Password      string                 `json:"password,omitempty"`
	PhoneNumber   string                 `json:"phone_number,omitempty"`
	UserMetadata  map[string]interface{} `json:"user_metadata,omitempty"`
	EmailVerified bool                   `json:"email_verified,omitempty"`
	VerifyEmail   bool                   `json:"verify_email,omitempty"`
	PhoneVerified bool                   `json:"phone_verified,omitempty"`
	AppMetadata   map[string]interface{} `json:"app_metadata,omitempty"`
	FamilyName    string                 `json:"family_name,omitempty"`
	GivenName     string                 `json:"given_name,omitempty"`
}

// UserUpdateOpts are options which can be used to update a user
type UserUpdateOpts struct {
	Blocked           bool                   `json:"blocked,omitempty"`
	EmailVerified     bool                   `json:"email_verified,omitempty"`
	Email             string                 `json:"email,omitempty"`
	VerifyEmail       bool                   `json:"verify_email,omitempty"`
	PhoneNumber       string                 `json:"phone_number,omitempty"`
	PhoneVerified     bool                   `json:"phone_verified,omitempty"`
	VerifyPhoneNumber bool                   `json:"verify_phone_number,omitempty"`
	Password          string                 `json:"password,omitempty"`
	VerifyPassword    bool                   `json:"verify_password,omitempty"`
	UserMetadata      map[string]interface{} `json:"user_metadata,omitempty"`
	AppMetadata       map[string]interface{} `json:"app_metadata,omitempty"`
	Connection        string                 `json:"connection,omitempty"`
	Username          string                 `json:"username,omitempty"`
	ClientID          string                 `json:"client_id,omitempty"`
}

// SearchUsersOpts are options which can be used to used to search users
type SearchUsersOpts struct {
	PerPage       int    `url:"per_page,omitempty"`
	Page          int    `url:"page,omitempty"`
	IncludeTotals bool   `url:"include_totals,omitempty"`
	Sort          string `url:"sort,omitempty"`
	Connection    string `url:"connection,omitempty"`
	Fields        string `url:"fields,omitempty"`
	IncludeFields bool   `url:"include_fields,omitempty"`
	Q             string `url:"q,omitempty"`
	SearchEngine  string `url:"search_engine,omitempty"`
}

// Encode creates a url.Values encoding of SearchUserOpts.
func (opts *SearchUsersOpts) Encode() (string, error) {
	vals, err := query.Values(opts)
	if err != nil {
		return "", err
	}
	return vals.Encode(), nil
}

// Identity is the identity of a user in Auth0
type Identity struct {
	Connection string `json:"connection,omitempty"`
	ID         string `json:"user_id,omitempty"`
	Provider   string `json:"provider,omitempty"`
	IsSocial   bool   `json:"isSocial,omitempty"`
}

// GetAll returns all users
func (svc *UsersService) GetAll() ([]User, error) {
	var users []User
	err := svc.c.GetV2("/users", &users)
	return users, err
}

// Get returns a users
func (svc *UsersService) Get(userID string) (User, error) {
	var user User
	err := svc.c.Get("/users/"+userID, &user)
	return user, err
}

// Search retrieves users according to search criteria
func (svc *UsersService) Search(opts SearchUsersOpts) (*UsersPage, error) {
    var usersPage UsersPage
	queryString, err := opts.Encode()
	if err != nil {
		return nil, err
	}
	url := "/users"
	if queryString != "" {
		url = fmt.Sprintf("/users?%s", queryString)
	}
    if opts.IncludeTotals {
	    err = svc.c.Get(url, &usersPage)
    } else {
        err = svc.c.Get(url, &usersPage.Users)
    }
	return &usersPage, err
}

// Create creates a user
func (svc *UsersService) Create(opts UserOpts) (User, error) {
	var user User
	err := svc.c.Post("/users", opts, &user)
	return user, err
}

// Delete deletes a users
func (svc *UsersService) Delete(userID string) error {
	return svc.c.Delete("/users/"+userID, nil, nil)
}

func (svc *UsersService) DeleteWithBody(userID string, body interface{}) error {
	return svc.c.Delete("/users/"+userID, &body, nil)
}

// Update updates a user
func (svc *UsersService) Update(userID string, opts UserUpdateOpts) (User, error) {
	var user User
	err := svc.c.Patch("/users/"+userID, &opts, &user)
	return user, err
}
