package omada

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var loginRedirectPattern = regexp.MustCompile(`/login$`)

// doRequest executes the request built by buildReq. If the response is a 302
// redirect to /login, it refreshes the session and retries once.
func (c *Controller) doRequest(buildReq func() (*http.Request, error)) (*http.Response, error) {
	req, err := buildReq()
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusFound {
		return res, nil
	}

	res.Body.Close()
	location := res.Header.Get("Location")
	if !loginRedirectPattern.MatchString(location) {
		return nil, fmt.Errorf("unexpected response: redirect to %s", location)
	}
	if err := c.refreshLogin(); err != nil {
		return nil, err
	}

	req, err = buildReq()
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

func invokeRequest[T any](c *Controller, path string, queryParams map[string]string) (*T, error) {

	address, err := url.JoinPath(c.baseURL, c.controllerId, path)
	if err != nil {
		return nil, err
	}
	omadaUrl, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	values := omadaUrl.Query()
	for k, v := range queryParams {
		values.Add(k, v)
	}
	omadaUrl.RawQuery = values.Encode()

	buildReq := func() (*http.Request, error) {
		req, err := http.NewRequest("GET", omadaUrl.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Add("Csrf-Token", c.token)
		return req, nil
	}

	res, err := c.doRequest(buildReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}

	var out T
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil

}

// invokePostRequest sends a POST to baseURL+path (controllerId is NOT prepended)
// with the given body marshalled as JSON and any extra headers.
func invokePostRequest[T any](c *Controller, path string, body any, extraHeaders map[string]string) (*T, error) {

	address, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, err
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	buildReq := func() (*http.Request, error) {
		req, err := http.NewRequest("POST", address, bytes.NewBuffer(bodyJSON))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Add("Csrf-Token", c.token)
		for k, v := range extraHeaders {
			req.Header.Set(k, v)
		}
		return req, nil
	}

	res, err := c.doRequest(buildReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}

	var out T
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil

}
