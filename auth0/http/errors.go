package http

import (
	"fmt"
)

// Error is an http error returned from the Auth0 service
type Error struct {
	StatusCode int    `json:"statusCode,omitempty"`
	HTTPError  string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (e Error) Error() string {
	msg := "auth0: "
	if e.StatusCode != 0 {
		msg += fmt.Sprintf("%d ", e.StatusCode)
	}

	if e.HTTPError != "" {
		msg += e.HTTPError + " "
	}

	if e.Message != "" {
		msg += "(" + e.Message + ")"
	}

	return msg
}
