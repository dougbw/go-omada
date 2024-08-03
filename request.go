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

	if res.StatusCode == http.StatusFound {
		fmt.Println("there is a login issue")
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		return nil, err
	}
	return res, nil

}
