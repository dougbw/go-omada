package omada

import (
	"encoding/json"
	"fmt"
	"sort"
)

type deviceResponse struct {
	ErrorCode int      `json:"errorCode"`
	Msg       string   `json:"msg"`
	Result    []Device `json:"result"`
}

type Device struct {
	Type            string `json:"type"`
	Mac             string `json:"mac"`
	Name            string `json:"name"`
	Model           string `json:"model"`
	CompoundModel   string `json:"compoundModel"`
	ShowModel       string `json:"showModel"`
	ModelVersion    string `json:"modelVersion"`
	FirmwareVersion string `json:"firmwareVersion"`
	Version         string `json:"version"`
	HwVersion       string `json:"hwVersion"`
	IP              string `json:"ip"`
	Uptime          string `json:"uptime"`
	UptimeLong      int    `json:"uptimeLong"`
	StatusCategory  int    `json:"statusCategory"`
	Status          int    `json:"status"`
	AdoptFailType   int    `json:"adoptFailType"`
	LastSeen        int64  `json:"lastSeen"`
	NeedUpgrade     bool   `json:"needUpgrade"`
	FwDownload      bool   `json:"fwDownload"`
	CPUUtil         int    `json:"cpuUtil"`
	MemUtil         int    `json:"memUtil"`
	Download        int64  `json:"download"`
	Upload          int64  `json:"upload"`
	Site            string `json:"site"`
	Location        struct {
		MapID       string  `json:"mapId"`
		PosX        float64 `json:"posX"`
		PosY        float64 `json:"posY"`
		InstallType int     `json:"installType"`
		Height      float64 `json:"height"`
		Located     bool    `json:"located"`
	} `json:"location,omitempty"`
	ClientNum       int    `json:"clientNum"`
	Compatible      int    `json:"compatible"`
	LocateEnable    bool   `json:"locateEnable"`
	Sn              string `json:"sn"`
	CombinedGateway bool   `json:"combinedGateway"`
	WirelessLinked  bool   `json:"wirelessLinked,omitempty"`
	DeviceMisc      struct {
		Support5G           bool `json:"support5g"`
		Support5G2          bool `json:"support5g2"`
		Support6G           bool `json:"support6g"`
		Support11Ac         bool `json:"support11ac"`
		SupportLag          bool `json:"supportLag"`
		SupportMesh         int  `json:"supportMesh"`
		CustomizeRegion     int  `json:"customizeRegion"`
		MinPower2G          int  `json:"minPower2G"`
		MaxPower2G          int  `json:"maxPower2G"`
		MinPower5G          int  `json:"minPower5G"`
		MaxPower5G          int  `json:"maxPower5G"`
		SupportChannelLimit bool `json:"supportChannelLimit"`
		SupportDfs          int  `json:"supportDfs"`
		SupportRoaming      int  `json:"supportRoaming"`
	} `json:"deviceMisc,omitempty"`
	DevCap struct {
		SupportPa           int  `json:"supportPa"`
		MeshChainNum        int  `json:"meshChainNum"`
		SupportOFDMA2G      bool `json:"supportOFDMA2g"`
		SupportOFDMA5G      bool `json:"supportOFDMA5g"`
		SupportOFDMA5G2     bool `json:"supportOFDMA5g2"`
		SupportOFDMA6G      bool `json:"supportOFDMA6g"`
		SupportMeshPriority bool `json:"supportMeshPriority"`
		SupportL3Access     bool `json:"supportL3Access"`
	} `json:"devCap,omitempty"`
	WlanGroup      string   `json:"wlanGroup,omitempty"`
	Override       string   `json:"override,omitempty"`
	Bssids         []string `json:"bssids,omitempty"`
	RadioSetting2G struct {
		RadioEnable  bool   `json:"radioEnable"`
		ChannelWidth string `json:"channelWidth"`
		Channel      string `json:"channel"`
		TxPower      int    `json:"txPower"`
		TxPowerLevel int    `json:"txPowerLevel"`
	} `json:"radioSetting2g,omitempty"`
	RadioSetting5G struct {
		RadioEnable  bool   `json:"radioEnable"`
		ChannelWidth string `json:"channelWidth"`
		Channel      string `json:"channel"`
		TxPower      int    `json:"txPower"`
		TxPowerLevel int    `json:"txPowerLevel"`
	} `json:"radioSetting5g,omitempty"`
	Wp2G struct {
		ActualChannel string `json:"actualChannel"`
		MaxTxRate     int    `json:"maxTxRate"`
		TxPower       int    `json:"txPower"`
		Region        int    `json:"region"`
		BandWidth     string `json:"bandWidth"`
		RdMode        string `json:"rdMode"`
		TxUtil        int    `json:"txUtil"`
		RxUtil        int    `json:"rxUtil"`
		InterUtil     int    `json:"interUtil"`
	} `json:"wp2g,omitempty"`
	Wp5G struct {
		ActualChannel string `json:"actualChannel"`
		MaxTxRate     int    `json:"maxTxRate"`
		TxPower       int    `json:"txPower"`
		Region        int    `json:"region"`
		BandWidth     string `json:"bandWidth"`
		RdMode        string `json:"rdMode"`
		TxUtil        int    `json:"txUtil"`
		RxUtil        int    `json:"rxUtil"`
		InterUtil     int    `json:"interUtil"`
	} `json:"wp5g,omitempty"`
	TxRate           int     `json:"txRate,omitempty"`
	RxRate           int     `json:"rxRate,omitempty"`
	ClientNum2G      int     `json:"clientNum2g,omitempty"`
	ClientNum5G      int     `json:"clientNum5g,omitempty"`
	ClientNum5G2     int     `json:"clientNum5g2,omitempty"`
	ClientNum6G      int     `json:"clientNum6g,omitempty"`
	UserNum          int     `json:"userNum,omitempty"`
	GuestNum         int     `json:"guestNum,omitempty"`
	Hop              int     `json:"hop,omitempty"`
	Downlink         int     `json:"downlink,omitempty"`
	AnyPoeEnable     bool    `json:"anyPoeEnable,omitempty"`
	LicenseStatusStr string  `json:"licenseStatusStr"`
	Uplink           string  `json:"uplink,omitempty"`
	LoopbackNum      int     `json:"loopbackNum,omitempty"`
	Loop             string  `json:"loop,omitempty"`
	PoeRemain        float64 `json:"poeRemain,omitempty"`
	FanStatus        int     `json:"fanStatus,omitempty"`
	PoeSupport       bool    `json:"poeSupport,omitempty"`
	DnsName          string
}

func (c *Controller) GetDevices() ([]Device, error) {

	path := fmt.Sprintf("api/v2/sites/%s/devices", c.siteId)
	queryParams := map[string]string{
		"currentPage":     "1",
		"currentPageSize": "999",
	}
	res, err := c.invokeRequest(path, queryParams)
	if err != nil {
		return nil, err
	}

	var deviceResponse deviceResponse
	if err := json.NewDecoder(res.Body).Decode(&deviceResponse); err != nil {
		return nil, err
	}

	if deviceResponse.ErrorCode != 0 {
		err = fmt.Errorf("failed to get list of devices: code='%d', message='%s'", deviceResponse.ErrorCode, deviceResponse.Msg)
		return nil, err
	}

	var devices []Device
	for _, device := range deviceResponse.Result {
		device.DnsName = makeDNSSafe(device.Name)
		devices = append(devices, device)
	}

	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Name < devices[j].Name
	})

	return devices, nil

}
