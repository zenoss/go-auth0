package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Doer can do http requests
type Doer interface {
	Do(*http.Request, interface{}) error
}

// Client handles requests to API
type Client struct {
	Doer
	API string
}

// RootClient is composed of an actual http.Client that makes the requests
type RootClient struct {
	*http.Client
}

func readAndUnmarshal(r io.Reader, obj interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "Cannot read response body")
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return errors.Wrap(err, "Cannot unmarshal response")
	}
	return nil
}

func getResponseError(resp *http.Response) error {
	if resp.ContentLength == 0 {
		return &Error{
			StatusCode: resp.StatusCode,
			HTTPError:  resp.Status,
		}
	}
	var respError Error
	defer resp.Body.Close()
	err := readAndUnmarshal(resp.Body, &respError)
	if err != nil {
		return err
	}
	return respError
}

// Do processes a request and unmarshals the response body into respBody
func (c *RootClient) Do(req *http.Request, respBody interface{}) error {
	// POSTs are application/json to this api
	if req.ContentLength > 0 && (req.Method == "POST" ||
		req.Method == "PUT" || req.Method == "PATCH") {
		req.Header.Add("Content-Type", "application/json")
	}
	// Perform the request
	resp, err := c.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Cannot complete request")
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// if we have a success code and no response body, we're done
		if resp.ContentLength == 0 {
			return nil
		}
		// if we have a response body, unmarshal it
		defer resp.Body.Close()
		return readAndUnmarshal(resp.Body, respBody)
	}
	return getResponseError(resp)
}

// Do processes a request and unmarshals the response body into respBody
func (c *Client) Do(req *http.Request, respBody interface{}) error {
	return c.Doer.Do(req, respBody)
}

func noSlash(uri string) string {
	return strings.TrimRight(uri, "/")
}

// Get performs a get to the endpoint of the API associated with the client
func (c *Client) GetWithHeaders(endpoint string, respBody interface{}, headers map[string]string) error {
	req, err := http.NewRequest("GET", noSlash(c.API)+endpoint, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}
	return c.Doer.Do(req, respBody)
}

// Get performs a get to the endpoint of the API associated with the client
func (c *Client) Get(endpoint string, respBody interface{}) error {
	return c.GetWithHeaders(endpoint, respBody, map[string]string{})
}

// Post performs a post to the endpoint of the API associated with the client
func (c *Client) PostWithHeaders(endpoint string, body interface{}, respBody interface{}, headers map[string]string) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("POST", noSlash(c.API)+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}
	return c.Doer.Do(req, respBody)
}

// Post performs a post to the endpoint of the API associated with the client
func (c *Client) Post(endpoint string, body interface{}, respBody interface{}) error {
	return c.PostWithHeaders(endpoint, body, respBody, map[string]string{})
}

// Put performs a put to the endpoint of the API associated with the client
func (c *Client) PutWithHeaders(endpoint string, body interface{}, respBody interface{}, headers map[string]string) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("PUT", noSlash(c.API)+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}
	return c.Doer.Do(req, respBody)
}

// Put performs a put to the endpoint of the API associated with the client
func (c *Client) Put(endpoint string, body interface{}, respBody interface{}) error {
	return c.PutWithHeaders(endpoint, body, respBody, map[string]string{})
}

// Patch performs a patch to the endpoint of the API associated with the client
func (c *Client) PatchWithHeaders(endpoint string, body interface{}, respBody interface{}, headers map[string]string) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("PATCH", noSlash(c.API)+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}
	return c.Doer.Do(req, respBody)
}

// Patch performs a patch to the endpoint of the API associated with the client
func (c *Client) Patch(endpoint string, body interface{}, respBody interface{}) error {
	return c.PatchWithHeaders(endpoint, body, respBody, map[string]string{})
}

// Delete performs a delete to the endpoint of the API associated with the client
func (c *Client) DeleteWithHeaders(endpoint string, body interface{}, respBody interface{}, headers map[string]string) error {
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Cannot marshal body")
	}
	req, err := http.NewRequest("DELETE", noSlash(c.API)+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}
	return c.Doer.Do(req, respBody)
}

// Delete performs a delete to the endpoint of the API associated with the client
func (c *Client) Delete(endpoint string, body interface{}, respBody interface{}) error {
	return c.DeleteWithHeaders(endpoint, body, respBody, map[string]string{})
}
