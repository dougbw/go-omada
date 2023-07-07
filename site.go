package omada

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Controller) getSites() error {

	path := "api/v2/users/current"
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, c.controllerId, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Csrf-Token", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		return err
	}

	var currentUserResponse currentUserResponse
	if err := json.NewDecoder(res.Body).Decode(&currentUserResponse); err != nil {
		return err
	}

	c.sites = make(map[string]string)
	for _, v := range currentUserResponse.Result.Privilege.Sites {
		c.sites[v.Name] = v.Key
		c.SetSite(v.Name)
	}

	return nil

}

func (c *Controller) SetSite(site string) error {
	siteId, ok := c.sites[site]
	if !ok {
		return fmt.Errorf("site not found: %s", site)
	}
	c.siteId = siteId
	return nil
}
