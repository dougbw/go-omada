package omada

import "fmt"

func (c *Controller) SetSite(site string) error {
	siteId, ok := c.sites[site]
	if !ok {
		return fmt.Errorf("site not found: %s", site)
	}
	c.siteId = siteId
	return nil
}
