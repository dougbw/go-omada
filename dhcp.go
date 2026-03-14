package omada

import (
	"fmt"
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

	path := fmt.Sprintf("api/v2/sites/%s/setting/service/dhcp", c.siteId)
	queryParams := map[string]string{
		"currentPage":     "1",
		"currentPageSize": "999",
	}
	res, err := invokeRequest[dhcpReservationResponse](c, path, queryParams)
	if err != nil {
		return nil, err
	}

	if res.ErrorCode != 0 {
		err = fmt.Errorf("failed to get list of dhcp reservation: code='%d', message='%s'", res.ErrorCode, res.Msg)
		return nil, err
	}

	return res.Result.Data, nil
}
