package omada

import (
	"fmt"
)

func (c *Controller) getSites() error {

	path := "api/v2/users/current"
	res, err := invokeRequest[currentUserResponse](c, path, nil)
	if err != nil {
		return err
	}

	if res.ErrorCode != 0 {
		return fmt.Errorf("failed to list sites: code='%d', message='%s'", res.ErrorCode, res.Msg)
	}

	c.Sites = make(map[string]string)
	for _, v := range res.Result.Privilege.Sites {
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
