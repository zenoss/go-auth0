package mgmt

import (
	"fmt"
	"github.com/zenoss/go-auth0/auth0/http"
	"net/url"
)

type TokenDataOpts struct {
	ID         string `json:"user_id,omitempty"`
	TokenType  string `json:"type,omitempty"`
}

type TokenData struct {
	DeviceName string `json:"device_name,omitempty"`
	ID         string `json:"id,omitempty"`
	TokenType  string `json:"type,omitempty"`
	UserID     string `json:"user_id,omitempty"`
}

type DeviceCredentials struct {
	c *http.Client
}

// Lists refresh tokens
func (svc *DeviceCredentials) Get(userID string) ([]TokenData, error) {
	//https://manage.auth0.com/api/device-credentials?user_id=auth0%7Ce8ey6zc9hfxppbz2h88r5yqqj&type=refresh_token
	var tokens []TokenData

	v := url.Values{}
	v.Set("user_id", userID)
	v.Add("type", "refresh_token")

	u := fmt.Sprintf("/device-credentials?%s", v.Encode())
	err := svc.c.GetV2(u, &tokens)
	return tokens, err
}

// Deletes all tokens for the user with a matching device identifier.
func (svc *DeviceCredentials) DeleteByIdentifier(userID, device string) error {
	tokens, err := svc.Get(userID)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		if token.DeviceName == device {
			err = svc.Delete(token.ID)
			if err != nil {
				// abort the loop; some tokens may be left intact in auth0. errors are typically a problem
				// with the call/scopes, so generally if one of these succeeds, they will all succeed.
				return err
			}
		}
	}
	return nil
}

// Deletes the specified token (this requires the auth0 id for the device, not the identifier)
func (svc *DeviceCredentials) Delete(tokenid string) error {
	endpoint := fmt.Sprintf("/device-credentials/"+tokenid)
	return svc.c.Delete(endpoint, nil, nil)
}
