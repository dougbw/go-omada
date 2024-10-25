package omada

import (
	"encoding/json"
	"fmt"
	"sort"
)

type clientResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    struct {
		TotalRows   int      `json:"totalRows"`
		CurrentPage int      `json:"currentPage"`
		CurrentSize int      `json:"currentSize"`
		Data        []Client `json:"data"`
	} `json:"result"`
}

type Client struct {
	Name     string `json:"name"`
	HostName string `json:"hostName,omitempty"`
	Ip       string `json:"ip"`
	MAC      string `json:"mac"`
	DnsName  string
}

func (c *Controller) GetClients() ([]Client, error) {

	path := fmt.Sprintf("api/v2/sites/%s/clients", c.siteId)
	queryParams := map[string]string{
		"currentPage":     "1",
		"currentPageSize": "999",
		"filters.active":  "true",
	}
	res, err := c.invokeRequest(path, queryParams)
	if err != nil {
		return nil, err
	}

	var clientResponse clientResponse
	if err := json.NewDecoder(res.Body).Decode(&clientResponse); err != nil {
		return nil, err
	}

	if clientResponse.ErrorCode != 0 {
		err = fmt.Errorf("failed to get list of clients: code='%d', message='%s'", clientResponse.ErrorCode, clientResponse.Msg)
		return nil, err
	}

	var clients []Client
	for _, client := range clientResponse.Result.Data {
		if client.Ip == "" {
			continue
		}

		var dnsName string
		dnsName = client.Name
		if client.Name == client.MAC && client.HostName != "--" {
			dnsName = client.HostName
		}
		client.DnsName = makeDNSSafe(dnsName)
		clients = append(clients, client)
	}

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].DnsName < clients[j].DnsName
	})

	return clients, nil

}
