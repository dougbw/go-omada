package omada

import (
	"fmt"
	"sort"
)

type GetNetworksResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    struct {
		TotalRows   int            `json:"totalRows"`
		CurrentPage int            `json:"currentPage"`
		CurrentSize int            `json:"currentSize"`
		Data        []OmadaNetwork `json:"data"`
	} `json:"result"`
}

type OmadaNetwork struct {
	Id      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Domain  string `json:"domain,omitempty"`
	Purpose string `json:"purpose"`
	Subnet  string `json:"gatewaySubnet"`
}

func (c *Controller) GetNetworks() ([]OmadaNetwork, error) {

	path := fmt.Sprintf("api/v2/sites/%s/setting/lan/networks", c.siteId)
	queryParams := map[string]string{
		"currentPage":     "1",
		"currentPageSize": "999",
	}
	res, err := invokeRequest[GetNetworksResponse](c, path, queryParams)
	if err != nil {
		return nil, err
	}

	if res.ErrorCode != 0 {
		err = fmt.Errorf("failed to get list of networks: code='%d', message='%s'", res.ErrorCode, res.Msg)
		return nil, err
	}

	networks := res.Result.Data
	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Name < networks[j].Name
	})

	return networks, nil

}
