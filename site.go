package omada

import (
	"encoding/json"
	"fmt"
)

func (c *Controller) getSites() error {

	path := "api/v2/users/current"
	res, err := c.invokeRequest(path, nil)
	if err != nil {
		return err
	}

	var currentUserResponse currentUserResponse
	if err := json.NewDecoder(res.Body).Decode(&currentUserResponse); err != nil {
		return err
	}

	if currentUserResponse.ErrorCode != 0 {
		return fmt.Errorf("failed to list sites: code='%d', message='%s'", currentUserResponse.ErrorCode, currentUserResponse.Msg)
	}

	c.Sites = make(map[string]string)
	for _, v := range currentUserResponse.Result.Privilege.Sites {
		c.Sites[v.Name] = v.Key
		c.SetSite(v.Name)
	}

	return nil

}

func (c *Controller) SetSite(site string) error {
	siteId, ok := c.Sites[site]
	if !ok {
		return fmt.Errorf("site not found: %s", site)
	}
	c.siteId = siteId
	return nil
}
