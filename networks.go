package omada

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

type GetNetworksResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    struct {
		TotalRows   int                  `json:"totalRows"`
		CurrentPage int                  `json:"currentPage"`
		CurrentSize int                  `json:"currentSize"`
		Data        []OmadaNetworkHelper `json:"data"`
	} `json:"result"`
}

type OmadaNetwork struct {
	Id         string `json:"id"`
	Name       string `json:"name,omitempty"`
	Domain     string `json:"domain,omitempty"`
	Purpose    string `json:"purpose"`
	Subnet     string `json:"gatewaySubnet"`
	Ipv6Subnet string
}

type OmadaNetworkHelper struct {
	Id                   string `json:"id"`
	Name                 string `json:"name,omitempty"`
	Domain               string `json:"domain,omitempty"`
	Purpose              string `json:"purpose"`
	Subnet               string `json:"gatewaySubnet"`
	LanNetworkIpv6Config struct {
		Proto  string `json:"proto"`
		Enable int    `json:"enable"`
		Slaac  struct {
			Prefix string `json:"prefix"`
		} `json:"slaac"`
		Rdnss struct {
			Prefix string `json:"prefix"`
		} `json:"rdnss"`
		Dhcpv6 struct {
			Gateway string `json:"gateway"`
			Subnet  int    `json:"subnet"`
		} `json:"dhcpv6"`
	} `json:"lanNetworkIpv6Config"`
}

func (c *Controller) GetNetworks() ([]OmadaNetwork, error) {

	path := fmt.Sprintf("api/v2/sites/%s/setting/lan/networks", c.siteId)
	queryParams := map[string]string{
		"currentPage":     "1",
		"currentPageSize": "999",
	}
	res, err := c.invokeRequest(path, queryParams)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var networkResponse GetNetworksResponse
	if err := json.NewDecoder(res.Body).Decode(&networkResponse); err != nil {
		return nil, err
	}

	if networkResponse.ErrorCode != 0 {
		err = fmt.Errorf("failed to get list of networks: code='%d', message='%s'", networkResponse.ErrorCode, networkResponse.Msg)
		return nil, err
	}

	networkhelpers := networkResponse.Result.Data

	networks := ConvertHelpersToNetworks(networkhelpers)

	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Name < networks[j].Name
	})

	return networks, nil

}

func ConvertHelpersToNetworks(networkHelpers []OmadaNetworkHelper) []OmadaNetwork {

	networks := make([]OmadaNetwork, 0, len(networkHelpers))

	for _, networkHelper := range networkHelpers {
		var network OmadaNetwork

		network.Id = networkHelper.Id
		network.Name = networkHelper.Name
		network.Domain = networkHelper.Domain
		network.Purpose = networkHelper.Purpose
		network.Subnet = networkHelper.Subnet

		if networkHelper.LanNetworkIpv6Config.Enable != 0 {
			switch networkHelper.LanNetworkIpv6Config.Proto {
			case "slaac":
				network.Ipv6Subnet = networkHelper.LanNetworkIpv6Config.Slaac.Prefix + "/64"
			case "rdnss":
				network.Ipv6Subnet = networkHelper.LanNetworkIpv6Config.Rdnss.Prefix + "/64"
			case "dhcpv6":
				network.Ipv6Subnet =
					networkHelper.LanNetworkIpv6Config.Dhcpv6.Gateway +
						"/" +
						strconv.Itoa(networkHelper.LanNetworkIpv6Config.Dhcpv6.Subnet)
			}
		}

		networks = append(networks, network)

	}

	return networks

}
