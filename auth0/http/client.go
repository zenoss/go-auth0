package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

// Doer can do http requests
type Doer interface {
	Do(*http.Request, any) error
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

func readAndUnmarshal(r io.Reader, obj any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("Cannot read response body: %w", err)
	}

	err = json.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal response: %w", err)
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
func (c *RootClient) Do(req *http.Request, respBody any) error {
	// POSTs are application/json to this api
	if req.ContentLength > 0 && (req.Method == http.MethodPost ||
		req.Method == http.MethodPut || req.Method == http.MethodPatch) {
		req.Header.Add("Content-Type", "application/json")
	}
	// Perform the request
	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Cannot complete request: %w", err)
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
func (c *Client) Do(req *http.Request, respBody any) error {
	return c.Doer.Do(req, respBody)
}

func noSlash(uri string) string {
	return strings.TrimRight(uri, "/")
}

// Get performs a get to the endpoint of the API associated with the client
func (c *Client) GetWithHeaders(endpoint string, respBody any, headers map[string]string) error {
	req, err := http.NewRequest(http.MethodGet, noSlash(c.API)+endpoint, http.NoBody)
	if err != nil {
		return fmt.Errorf("Cannot create request: %w", err)
	}

	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}

	return c.Doer.Do(req, respBody)
}

// Get performs a get to the endpoint of the API associated with the client
func (c *Client) Get(endpoint string, respBody any) error {
	return c.GetWithHeaders(endpoint, respBody, map[string]string{})
}

// Get performs a get to the endpoint of the API v2 associated with the client
//
//revive:disable:cognitive-complexity
func (c *Client) GetWithHeadersV2(endpoint string, respBody any, headers map[string]string) error {
	// Support for a previous version of auth0 api
	fullUrl := noSlash(c.API) + endpoint
	if !strings.HasSuffix(c.API, "v2") {
		response, err := makeGetRequest(fullUrl, headers, c.Doer.Do)
		if err != nil {
			return err
		}

		return convertResponseData(response, respBody)
	}

	// auth0 v2 api returns maxPage 100 elements per page
	maxPage := 100
	page := 0
	fullUrl = addPagingParams(fullUrl, page, maxPage)
	keyName := extractKeyFromEndpoint(fullUrl)

	response, err := makeGetRequest(fullUrl, headers, c.Doer.Do)
	if err != nil {
		return err
	}

	var results []any

	var total int

	if val, ok := response.(map[string]any); ok {
		if t, ok := val["total"]; ok {
			total = int(t.(float64))
		}

		if items, ok := val[keyName]; ok {
			results = append(results, items.([]any)...)
		}
	}

	if total <= maxPage {
		return convertResponseData(results, respBody)
	}

	chanLen := (total / maxPage) + 1
	data := make(chan any, chanLen)
	g := errgroup.Group{}

	// spawn a bounded number of goroutines
	urls := make(chan string, chanLen)

	for i := maxPage; i < total; i += maxPage {
		page += 1
		// queue up the requests
		urls <- addPagingParams(fullUrl, page, maxPage)
	}

	close(urls)

	limiter := rate.NewLimiter(2, 2)

	for range 2 {
		g.Go(func() error {
			for fullUrl := range urls {
				_ = limiter.Wait(context.TODO())

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
		if val, ok := d.(map[string]any); ok {
			if items, ok := val[keyName]; ok {
				results = append(results, items.([]any)...)
			}
		}
	}

	return convertResponseData(results, respBody)
}

//revive:enable:cognitive-complexity

// Get performs a get to the endpoint of the API v2 associated with the client
func (c *Client) GetV2(endpoint string, respBody any) error {
	return c.GetWithHeadersV2(endpoint, respBody, map[string]string{})
}

// Get performs a get to the endpoint of the API v2 associated with the client,
// only for the summary, and returns the record count.
func (c *Client) CountWithHeadersV2(endpoint string, headers map[string]string) (int, error) {
	fullUrl := noSlash(c.API) + endpoint
	fullUrl = addPagingParams(fullUrl, 0, 1)

	response, err := makeGetRequest(fullUrl, headers, c.Doer.Do)
	if err != nil {
		return 0, err
	}

	if val, ok := response.(map[string]any); ok {
		if t, ok := val["total"]; ok {
			return int(t.(float64)), nil
		}

		return 0, fmt.Errorf("No total record count returned by GET %s query", fullUrl)
	}

	return 0, fmt.Errorf("Unable to process response to GET %s query", fullUrl)
}

// Get performs a get to the endpoint of the API v2 associated with the client,
// but returns the number of records rather than the actual data.
func (c *Client) CountV2(endpoint string) (int, error) {
	return c.CountWithHeadersV2(endpoint, map[string]string{})
}

// Post performs a post to the endpoint of the API associated with the client
func (c *Client) PostWithHeaders(endpoint string, body any, respBody any, headers map[string]string) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Cannot marshal body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, noSlash(c.API)+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Cannot create request: %w", err)
	}

	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}

	return c.Doer.Do(req, respBody)
}

// Post performs a post to the endpoint of the API associated with the client
func (c *Client) Post(endpoint string, body any, respBody any) error {
	return c.PostWithHeaders(endpoint, body, respBody, map[string]string{})
}

// Put performs a put to the endpoint of the API associated with the client
func (c *Client) PutWithHeaders(endpoint string, body any, respBody any, headers map[string]string) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Cannot marshal body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, noSlash(c.API)+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Cannot create request: %w", err)
	}

	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}

	return c.Doer.Do(req, respBody)
}

// Put performs a put to the endpoint of the API associated with the client
func (c *Client) Put(endpoint string, body any, respBody any) error {
	return c.PutWithHeaders(endpoint, body, respBody, map[string]string{})
}

// Patch performs a patch to the endpoint of the API associated with the client
func (c *Client) PatchWithHeaders(endpoint string, body any, respBody any, headers map[string]string) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Cannot marshal body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, noSlash(c.API)+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Cannot create request: %w", err)
	}

	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}

	return c.Doer.Do(req, respBody)
}

// Patch performs a patch to the endpoint of the API associated with the client
func (c *Client) Patch(endpoint string, body any, respBody any) error {
	return c.PatchWithHeaders(endpoint, body, respBody, map[string]string{})
}

// Delete performs a delete to the endpoint of the API associated with the client
func (c *Client) DeleteWithHeaders(endpoint string, body any, respBody any, headers map[string]string) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Cannot marshal body: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, noSlash(c.API)+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Cannot create request: %w", err)
	}

	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}

	return c.Doer.Do(req, respBody)
}

// Delete performs a delete to the endpoint of the API associated with the client
func (c *Client) Delete(endpoint string, body any, respBody any) error {
	return c.DeleteWithHeaders(endpoint, body, respBody, map[string]string{})
}

func extractKeyFromEndpoint(fullUrl string) string {
	// endpoint can be equal to "/users" or "/device-credentials?user_id=%s&type=refresh_token"
	u, _ := url.Parse(fullUrl)
	path := u.Path
	li := strings.LastIndex(path, "/")
	key := path[li+1:]

	return strings.ReplaceAll(key, "-", "_")
}

func convertResponseData(data any, container any) error {
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
	values.Set("page", strconv.Itoa(page))
	values.Set("per_page", strconv.Itoa(perPage))
	values.Set("include_totals", "true")
	u.RawQuery = values.Encode()

	return u.String()
}

func makeGetRequest(fullUrl string, headers map[string]string, requester func(*http.Request, any) error) (any, error) {
	req, err := http.NewRequest(http.MethodGet, fullUrl, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("Cannot create request: %w", err)
	}

	for key, value := range headers {
		if len(strings.TrimSpace(key)) > 0 && len(strings.TrimSpace(value)) > 0 {
			req.Header.Add(key, value)
		}
	}

	var temporaryResponse any

	err = requester(req, &temporaryResponse)
	if err != nil {
		return nil, err
	}

	return temporaryResponse, nil
}
