package omada

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClients(t *testing.T) {

	controllerId := "123bee230c77bbb45d9c8545d04d700a"
	siteId := "Default"
	pathLogin := fmt.Sprintf("/%s/api/v2/login", controllerId)
	pathUsers := fmt.Sprintf("/%s/api/v2/users/current", controllerId)
	pathClients := fmt.Sprintf("/%s/api/v2/sites/%s/clients", controllerId, siteId)
	pathDevices := fmt.Sprintf("/%s/api/v2/sites/%s/devices", controllerId, siteId)
	pathNetworks := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/networks", controllerId, siteId)

	responses := map[string]string{
		"/api/info":  "./test-data/info-response.json",
		pathLogin:    "./test-data/info-response.json",
		pathUsers:    "./test-data/users-response.json",
		pathClients:  "./test-data/clients-response.json",
		pathDevices:  "./test-data/devices-response.json",
		pathNetworks: "./test-data/networks-response.json",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseFile, ok := responses[r.URL.Path]
		if !ok {
			t.Errorf("Unexpected request path on mock server: %s", r.URL.Path)
		}
		response, err := os.ReadFile(responseFile)
		if err != nil {
			fmt.Print(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}))
	defer server.Close()

	controller := New(server.URL)
	err := controller.GetControllerInfo()
	if err != nil {
		t.Fatalf("test failure on 'GetControllerInfo': %v", err)
	}
	err = controller.Login("user", "pass")
	if err != nil {
		t.Fatalf("test failure on 'Login': %v", err)
	}

	err = controller.SetSite("Home")
	if err != nil {
		t.Fatalf("test failure on 'SetSite': %v", err)
	}

	clients, err := controller.GetClients()
	if err != nil {
		t.Fatalf("test failure on 'GetClients': %v", err)
	}

	assert.Equal(t, controller.siteId, "Default")
	expectedClients := 3
	assert.Len(t, clients, expectedClients)
	assert.Equal(t, "Client 001", clients[0].Name)
	assert.Equal(t, "10.0.0.101", clients[0].Ip)
	assert.Equal(t, "client-001", clients[0].DnsName)

	assert.Equal(t, "Google Nest Mini", clients[1].Name)
	assert.Equal(t, "10.0.0.103", clients[1].Ip)
	assert.Equal(t, "google-nest-mini", clients[1].DnsName)

}
