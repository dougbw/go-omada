package omada

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Controller) sendRequest(req *http.Request) (*http.Response, error) {

	type responseBody struct {
		ErrorCode int    `json:"errorCode"`
		Msg       string `json:"msg"`
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		return nil, err
	}

	var buf bytes.Buffer
	var response responseBody
	tee := io.TeeReader(res.Body, &buf)
	if err := json.NewDecoder(tee).Decode(&response); err != nil {
		return nil, err
	}

	if response.ErrorCode != 0 {
		fmt.Printf("response error code: %d", response.ErrorCode)

		// attempt login
		err := c.refreshLogin()
		if err != nil {
			return nil, err
		}

		// retry
		res, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("status code: %d", res.StatusCode)
			return nil, err
		}

	}

	return res, nil
}

func (c *Controller) invokeRequest(path string, queryParams map[string]string) (*http.Response, error) {

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

	req, err := http.NewRequest("GET", omadaUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Csrf-Token", c.token)

	res, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return res, nil

}
