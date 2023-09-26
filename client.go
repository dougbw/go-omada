package omada

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	token := c.token
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/clients?currentPage=1&currentPageSize=999&filters.active=true", c.baseURL, c.controllerId, c.siteId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Csrf-Token", token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		return nil, err
	}

	var clientResponse clientResponse
	if err := json.NewDecoder(res.Body).Decode(&clientResponse); err != nil {
		return nil, err
	}

	if clientResponse.ErrorCode != 0 {
		err = fmt.Errorf("failed to get list of clients. API response is: errorCode: '%d', message: '%s'", clientResponse.ErrorCode, clientResponse.Msg)
		return nil, err
	}

	var clients []Client
	for _, client := range clientResponse.Result.Data {
		if client.Ip == "" {
			continue
		}
		client.DnsName = makeDNSSafe(client.Name)
		clients = append(clients, client)
	}

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].DnsName < clients[j].DnsName
	})

	return clients, nil

}
