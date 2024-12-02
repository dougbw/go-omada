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

	Active                    bool     `json:"active"`
	Activity                  int64    `json:"activity"`
	APMac                     string   `json:"apMac,omitempty"`
	APName                    string   `json:"apName,omitempty"`
	AuthStatus                int      `json:"authStatus"`
	Channel                   int      `json:"channel,omitempty"`
	ConnectDevType            string   `json:"connectDevType,omitempty"`
	ConnectedToWirelessRouter bool     `json:"connectedToWirelessRouter"`
	ConnectType               int      `json:"connectType,omitempty"`
	DeviceType                string   `json:"deviceType,omitempty"`
	Dot1xVlan                 int      `json:"dot1xVlan"`
	DownPacket                int64    `json:"downPacket"`
	Guest                     bool     `json:"guest"`
	HealthScore               int      `json:"healthScore"`
	IPV6List                  []string `json:"ipv6List"`
	LastSeen                  int64    `json:"lastSeen"`
	Manager                   bool     `json:"manager"`
	NetworkName               string   `json:"networkName,omitempty"`
	Port                      int      `json:"port"`
	PowerSave                 bool     `json:"powerSave"`
	RadioID                   int      `json:"radioId,omitempty"`
	RSSI                      int      `json:"rssi"`
	RXRate                    int64    `json:"rxRate,omitempty"`
	SignalLevel               int      `json:"signalLevel"`
	SignalRank                int      `json:"signalRank"`
	SNR                       int      `json:"snr"`
	SSID                      string   `json:"ssid,omitempty"`
	StackableSwitch           bool     `json:"stackableSwitch"`
	StandardPort              string   `json:"standardPort,omitempty"`
	Support5G2                bool     `json:"support5g2"`
	SwitchMac                 string   `json:"switchMac,omitempty"`
	SwitchName                string   `json:"switchName,omitempty"`
	TrafficDown               int64    `json:"trafficDown"`
	TrafficUp                 int64    `json:"trafficUp"`
	TXRate                    int64    `json:"txRate,omitempty"`
	UpPacket                  int64    `json:"upPacket"`
	Uptime                    int64    `json:"uptime"`
	VID                       int      `json:"vid"`
	WifiMode                  int      `json:"wifiMode,omitempty"`
	Wireless                  bool     `json:"wireless"`
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
		client.DnsName = makeDNSSafe(client.Name)
		clients = append(clients, client)
	}

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].DnsName < clients[j].DnsName
	})

	return clients, nil

}
