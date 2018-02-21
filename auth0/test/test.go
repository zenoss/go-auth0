package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/zenoss/go-auth0/auth0"
	"github.com/zenoss/go-auth0/auth0/authz"
	authhttp "github.com/zenoss/go-auth0/auth0/http"
	"github.com/zenoss/go-auth0/auth0/mgmt"
)

const (
	AuthzRoute              = "/authz"
	AuthzGroupsRoute        = AuthzRoute + "/api/groups"
	AuthzGroupRoute         = AuthzGroupsRoute + "/{groupid}"
	AuthzGroupMappingsRoute = AuthzGroupRoute + "/mappings"
	AuthzGroupRolesRoute    = AuthzGroupRoute + "/roles"
)

// A ResponseFunc should create a response given the request
type ResponseFunc func(*http.Request) interface{}

// A RouteConfig defines configuration for adding a route to the test server
type RouteConfig struct {
	Route        string
	Response     interface{}
	ResponseFunc *ResponseFunc
	Methods      []string
	Headers      map[string]string
	Queries      map[string]string
}

type Auth0TestServer struct {
	*httptest.Server
	Handler *mux.Router
}

func NewAuth0TestServer() *Auth0TestServer {
	r := mux.NewRouter()
	s := &Auth0TestServer{
		Server:  httptest.NewServer(r),
		Handler: r,
	}
	return s
}

func (s *Auth0TestServer) Client() *auth0.Auth0 {
	c := &auth0.Auth0{
		Client: &authhttp.Client{
			Doer: &authhttp.RootClient{
				Client: s.Server.Client(),
			},
			Site: s.Server.URL,
		},
	}
	c.Token = &auth0.TokenService{
		Client: c.Client,
	}
	c.Mgmt = mgmt.New(c.Site, c.Client)
	c.Authz = authz.New(s.Server.URL+AuthzRoute, c.Client)
	return c
}

func (s *Auth0TestServer) AddRoute(cfg RouteConfig) {
	if cfg.Route == "" {
		return
	}
	r := s.Handler.Handle(cfg.Route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		var err error
		var response interface{}
		if cfg.ResponseFunc != nil {
			response = (*cfg.ResponseFunc)(r)
		} else {
			response = cfg.Response
		}
		if response != nil {
			data, err = json.Marshal(&response)
			if err != nil {
				fmt.Fprint(w, err.Error())
			} else {
				fmt.Fprint(w, string(data))
			}
		}
	}))
	if len(cfg.Methods) != 0 {
		r.Methods(cfg.Methods...)
	}
	for k, v := range cfg.Headers {
		r.Headers(k, v)
	}
	for k, v := range cfg.Queries {
		r.Queries(k, v)
	}
}

func (s *Auth0TestServer) AddRouteResponse(route string, response interface{}) {
	s.Handler.Handle(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		var err error
		if response != nil {
			data, err = json.Marshal(&response)
			if err != nil {
				fmt.Fprint(w, err.Error())
			} else {
				fmt.Fprint(w, string(data))
			}
		}
	}))
}

func (s *Auth0TestServer) AddAuthzRouteResponse(route string, response interface{}) {
	s.AddRouteResponse(s.Server.URL+AuthzRoute+route, response)
}
