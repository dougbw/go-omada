package omada

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type dhcpReservationResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    struct {
		TotalRows   int               `json:"totalRows"`
		CurrentPage int               `json:"currentPage"`
		CurrentSize int               `json:"currentSize"`
		Data        []DhcpReservation `json:"data"`
		LimitEnable bool              `json:"limitEnable"`
	} `json:"result"`
}

type DhcpReservation struct {
	ID                   string `json:"id"`
	Description          string `json:"description,omitempty"`
	NetID                string `json:"netId"`
	Mac                  string `json:"mac"`
	IP                   string `json:"ip"`
	Status               bool   `json:"status"`
	NetName              string `json:"netName"`
	ExportToIPMacBinding bool   `json:"exportToIpMacBinding"`
	ClientName           string `json:"clientName"`
}

func (c *Controller) GetDhcpReservations() ([]DhcpReservation, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/service/dhcp?currentPage=1&currentPageSize=999", c.baseURL, c.controllerId, c.siteId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Csrf-Token", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		return nil, err
	}

	var dhcpReservationResponse dhcpReservationResponse
	if err := json.NewDecoder(res.Body).Decode(&dhcpReservationResponse); err != nil {
		return nil, err
	}

	return dhcpReservationResponse.Result.Data, nil
}
