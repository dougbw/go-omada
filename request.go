package omada

import (
	"fmt"
	"net/http"
	"net/url"
)

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

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// if the response code is 302 then there might be a login issue
	// attempt to refresh the login and retry the request
	if res.StatusCode == http.StatusFound {
		err := c.refreshLogin()
		if err != nil {
			return nil, err
		}
		req.Header.Set("Csrf-Token", c.token)
		res2, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		if res2.StatusCode != http.StatusOK {
			err = fmt.Errorf("status code: %d", res.StatusCode)
			return nil, err
		}
		return res2, nil
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		return nil, err
	}
	return res, nil

}
