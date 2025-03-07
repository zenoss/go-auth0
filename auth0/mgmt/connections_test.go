//go:build integration
// +build integration

package mgmt_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zenoss/go-auth0/auth0/mgmt"
)

const (
	connectionName = "testconnection"
)

func createConnection(suite *ManagementTestSuite) mgmt.Connection {
	connection, err := suite.management.Connections.Create(mgmt.ConnectionOpts{
		Name:        connectionName,
		DisplayName: "Test Connection",
		Strategy:    "samlp",
		Options: map[string]interface{}{
			"signInEndpoint": "https://dev-55788114.okta.com/app/dev-55788114_zenoss_1/exk16k08yeoVd7v3y5d7/sso/saml",
			"signingCert":    "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURxRENDQXBDZ0F3SUJBZ0lHQVhwZVpDeXFNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1JR1VNUXN3Q1FZRFZRUUdFd0pWVXpFVE1CRUcNCkExVUVDQXdLUTJGc2FXWnZjbTVwWVRFV01CUUdBMVVFQnd3TlUyRnVJRVp5WVc1amFYTmpiekVOTUFzR0ExVUVDZ3dFVDJ0MFlURVUNCk1CSUdBMVVFQ3d3TFUxTlBVSEp2ZG1sa1pYSXhGVEFUQmdOVkJBTU1ER1JsZGkwMU5UYzRPREV4TkRFY01Cb0dDU3FHU0liM0RRRUoNCkFSWU5hVzVtYjBCdmEzUmhMbU52YlRBZUZ3MHlNVEEyTXpBeE9USTJOREZhRncwek1UQTJNekF4T1RJM05ERmFNSUdVTVFzd0NRWUQNClZRUUdFd0pWVXpFVE1CRUdBMVVFQ0F3S1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ3d05VMkZ1SUVaeVlXNWphWE5qYnpFTk1Bc0cNCkExVUVDZ3dFVDJ0MFlURVVNQklHQTFVRUN3d0xVMU5QVUhKdmRtbGtaWEl4RlRBVEJnTlZCQU1NREdSbGRpMDFOVGM0T0RFeE5ERWMNCk1Cb0dDU3FHU0liM0RRRUpBUllOYVc1bWIwQnZhM1JoTG1OdmJUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0MNCmdnRUJBSmVub05haER0bmFYREZISlZWZjFwUlAxbWh1UnFJemRScUZ5Rm11L0krVUg0YzlTZjMySDNXTGRmSDBqbDl1M0NHMWl0SGENCkNoekJHY0hCd1RxREJic1ZpdWFXaHZleTVNS1JmWllERnRHRmFQY2V5OHo1RVBMVG9WRndEOWZXMWhnSThqUjVPeUp0SnlGQ0o0TXoNCmFIbFVYZjhwVjJaVndmZEkxWkx4WHczSkVWTWRrTVEvUGFkSDlHcm5jUDhtSUdpdWhQbnRzZmlGTzAvK2MxdzN3ZFFSWkVxSXVZY2gNCjF3ekNaTXFzNkNjc1k4cWtLN3dEWDJWa2U2M0FOcjFuTnZQTGNnbUJ3cCtsWEtsLzJ4Q1FxQkJZUlhNUW5Ia2tRMzlzS0VMT21ZNXkNCnA4eDQ3M2U1cUNJUWV2TmhRb3ZlaGIvUS9TOFJneEI4YnVWaThOdEpSREVDQXdFQUFUQU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUENCkp6TEVTbEZBck9rcVNZVFRpdHVacFZaU2FrNk4vam5pYnBaSk44Q21KK3hDTnNwNDY3MFFVZUxwcDZ6bTZWcUxZV0F0cUFsR2VkZ24NCkNLVXBYb1o2aXhHa2tiS2RYSGdDbUVpcllDZ3ZZZE1vd3FadGxIWTgxbnpnT1NKS3pDNDU2Nk5ZNmxLZ0tMcjkrck5NZVFTdndJcTANCkpKVEpEZmlSUHdDYnFvNThHemxBNTByRlQ1ZmpkeGNDQWplTHBHUnViSUVtbFRBM1FaYXVtU3NLVjdnWnA4Y2VEanNURXoxMlF1aWMNCndoVmU2bkNGTXFtZ29TL2xMb0VRbyt4Nnl1cjdlRFJwa29vQ3pKc2lIK3JjbTJZVEE1WnFBUk81WUZqQm95TzJkSmpWa2MvOUFoUUwNCjFrM0dnSkU3U2ZSWEZpOVE2d2c0dzBPU1hBeHRwYzFSU25DOHpRPT0NCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
		},
	})

	assert.Nil(suite.T(), err)
	return connection
}

func getAllConnections(suite *ManagementTestSuite) []mgmt.Connection {
	connections, err := suite.management.Connections.GetAll()
	assert.Nil(suite.T(), err)
	return connections
}

func searchConnections(suite *ManagementTestSuite, searchOpts mgmt.SearchConnectionsOpts) *mgmt.ConnectionsPage {
	connectionsPage, err := suite.management.Connections.Search(searchOpts)
	assert.Nil(suite.T(), err)
	return connectionsPage
}

func deleteConnection(suite *ManagementTestSuite, ID string, ignoreErr bool) {
	err := suite.management.Connections.Delete(ID)
	if !ignoreErr {
		assert.Nil(suite.T(), err)
	}
}

func cleanUpConnections(suite *ManagementTestSuite) {
	connections := getAllConnections(suite)
	var ID string
	for _, connection := range connections {
		if connection.Name == connectionName {
			ID = connection.ID
			break
		}
	}
	if ID != "" {
		deleteConnection(suite, ID, true)
	}
}

func (suite *ManagementTestSuite) TestConnectionsCreateGetAllDelete() {
	t := suite.T()
	svc := suite.management.Connections

	// Check if connection existed before the test and remove
	cleanUpConnections(suite)

	// Create a connection
	connection := createConnection(suite)
	assert.Equal(t, connectionName, connection.Name)

	// Check we made it successfully
	connection, err := svc.Get(connection.ID)
	assert.Nil(t, err)
	assert.Equal(t, connectionName, connection.Name)

	// Check that we can search
	searchOpts := mgmt.SearchConnectionsOpts{
		Name: connectionName,
	}
	connections, err := svc.Search(searchOpts)
	assert.Nil(t, err)
	assert.NotNil(t, connections)

	// Update it
	update := mgmt.ConnectionUpdateOpts{
		DisplayName: "Test Connection Renamed",
	}
	connection, err = svc.Update(connection.ID, update)
	assert.Nil(t, err)
	assert.Equal(t, "Test Connection Renamed", connection.DisplayName)

	// Delete it
	deleteConnection(suite, connection.ID, true)

	// Check it was deleted
	_, err = svc.Get(connection.ID)
	assert.NotNil(t, err)
}
