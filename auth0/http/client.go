package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	defer func() {
		_ = resp.Body.Close()
	}()
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
	// get the rate limiter for this request
	limiter := GetRequestLimiter(req)
	reservation := limiter.Reserve()
	defer reservation.Cancel()
	if !reservation.OK() {
		time.Sleep(1200 * time.Millisecond)
	} else {
		reservation.Delay()
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
		defer func() {
			_ = resp.Body.Close()
		}()
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

// Get performs a get to the endpoint of the API v2 associated with the client
func (c *Client) GetWithHeadersV2(endpoint string, respBody interface{}, headers map[string]string) error {
	// Support for a previous version of auth0 api
	fullUrl := noSlash(c.API) + endpoint
	if !strings.HasSuffix(c.API, "v2") {
		response, err := makeGetRequest(fullUrl, headers, c.Doer.Do)
		if err != nil {
			return err
		}
		return convertResponseData(response, respBody)
	}

	//auth0 v2 api returns max 50 elements per page
	max := 50
	page := 0
	fullUrl = addPagingParams(fullUrl, page, max)
	keyName := extractKeyFromEndpoint(fullUrl)

	response, err := makeGetRequest(fullUrl, headers, c.Doer.Do)
	if err != nil {
		return err
	}

	var results []interface{}

	var total int
	if val, ok := response.(map[string]interface{}); ok {
		if t, ok := val["total"]; ok {
			total = int(t.(float64))
		}

		if items, ok := val[keyName]; ok {
			results = append(results, items.([]interface{})...)
		}
	}

	if total <= max {
		return convertResponseData(results, respBody)
	}

	chanLen := (total / max) + 1
	data := make(chan interface{}, chanLen)
	g := errgroup.Group{}

	// spawn a bounded number of goroutines
	urls := make(chan string)
	for i := max; i < total; i += max {
		page += 1
		// queue up the requests
		urls <- addPagingParams(fullUrl, page, max)
	}
	close(urls)
	for i := 0; i < 2; i++ {
		g.Go(func() error {
			for fullUrl := range urls {
				response, err := makeGetRequest(fullUrl, headers, c.Doer.Do)
				if err != nil {
					return err
				}
				data <- response
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	close(data)

	for d := range data {
		if val, ok := d.(map[string]interface{}); ok {
			if items, ok := val[keyName]; ok {
				results = append(results, items.([]interface{})...)
			}
		}
	}

	return convertResponseData(results, respBody)
}

// Get performs a get to the endpoint of the API v2 associated with the client
func (c *Client) GetV2(endpoint string, respBody interface{}) error {
	return c.GetWithHeadersV2(endpoint, respBody, map[string]string{})
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

func extractKeyFromEndpoint(fullUrl string) string {
	// endpoint can be equal to "/users" or "/device-credentials?user_id=%s&type=refresh_token"
	u, _ := url.Parse(fullUrl)
	path := u.Path
	li := strings.LastIndex(path, "/")
	key := path[li+1:]
	return strings.Replace(key, "-", "_", -1)
}

func convertResponseData(data interface{}, container interface{}) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataJson, &container)
	if err != nil {
		return err
	}
	return nil
}

func addPagingParams(fullUrl string, page, perPage int) string {
	u, _ := url.Parse(fullUrl)
	values, _ := url.ParseQuery(u.RawQuery)
	values.Set("page", fmt.Sprintf("%d", page))
	values.Set("per_page", fmt.Sprintf("%d", perPage))
	values.Set("include_totals", "true")
	u.RawQuery = values.Encode()

	return u.String()
}

func makeGetRequest(fullUrl string, headers map[string]string, requester func(*http.Request, interface{}) error) (interface{}, error) {
	req, err := http.NewRequest("GET", fullUrl, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot create request")
	}

	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}

	var temporaryResponse interface{}
	err = requester(req, &temporaryResponse)
	if err != nil {
		return nil, err
	}
	return temporaryResponse, nil
}
