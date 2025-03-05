package mgmt

import (
	"net/url"

	"github.com/zenoss/go-auth0/auth0/http"
)

type TokenDataOpts struct {
	ID        string `json:"user_id,omitempty"`
	TokenType string `json:"type,omitempty"`
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
	// https://manage.auth0.com/api/device-credentials?user_id=auth0%7Ce8ey6zc9hfxppbz2h88r5yqqj&type=refresh_token
	var tokens []TokenData

	v := url.Values{}
	v.Set("user_id", userID)
	v.Add("type", "refresh_token")
	u := "/device-credentials?" + v.Encode()
	err := svc.c.GetV2(u, &tokens)
	return tokens, err
}

// Count refresh tokens
func (svc *DeviceCredentials) Count(userID string) (int, error) {
	// https://manage.auth0.com/api/device-credentials?user_id=auth0%7Ce8ey6zc9hfxppbz2h88r5yqqj&type=refresh_token

	v := url.Values{}
	v.Set("user_id", userID)
	v.Add("type", "refresh_token")
	u := "/device-credentials?" + v.Encode()
	count, err := svc.c.CountV2(u)
	return count, err
}

// Deletes all tokens for the user with a matching device identifier.
func (svc *DeviceCredentials) DeleteByIdentifierInTokens(_, device string, tokens []TokenData) error {
	seenTokenId := map[string]bool{}
	for _, token := range tokens {
		if _, seen := seenTokenId[token.ID]; seen {
			// Ignore this duplicate token ID, most likely due to an auth0
			// bug we reported in https://support.auth0.com/tickets/00496600
			continue
		}
		seenTokenId[token.ID] = true

		if token.DeviceName == device {
			err := svc.Delete(token.ID)
			if err != nil {
				// abort the loop; some tokens may be left intact in auth0. errors are typically a problem
				// with the call/scopes, so generally if one of these succeeds, they will all succeed.
				return err
			}
		}
	}
	return nil
}

// Deletes all tokens for the user with a matching device identifier.
func (svc *DeviceCredentials) DeleteByIdentifier(userID, device string) error {
	tokens, err := svc.Get(userID)
	if err != nil {
		return err
	}

	return svc.DeleteByIdentifierInTokens(userID, device, tokens)
}

// Deletes the specified token (this requires the auth0 id for the device, not the identifier)
func (svc *DeviceCredentials) Delete(tokenid string) error {
	return svc.c.Delete("/device-credentials/"+tokenid, nil, nil)
}

// Deletes all grants for a user, which removes their refresh tokens across all devices.
func (svc *DeviceCredentials) DeleteGrants(userID string) error {
	v := url.Values{}
	v.Set("user_id", userID)
	return svc.c.Delete("/grants?"+v.Encode(), nil, nil)
}
