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

	controllerId := "978bee230c77bbb45d9c8545d04d700a"
	siteId := "Default"

	pathLogin := fmt.Sprintf("/%s/api/v2/login", controllerId)
	pathUsers := fmt.Sprintf("/%s/api/v2/users/current", controllerId)
	pathClients := fmt.Sprintf("/%s/api/v2/sites/%s/clients", controllerId, siteId)

	responseMap := map[string]string{
		"/api/info": "./test-data/info-response.json",
		pathLogin:   "./test-data/info-response.json",
		pathUsers:   "./test-data/users-response.json",
		pathClients: "./test-data/client-response.json",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Path {
		case "/api/info":
			response, err := os.ReadFile("./test-data/info-response.json")
			if err != nil {
				fmt.Print(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		case pathLogin:
			response, err := os.ReadFile("./test-data/login-response.json")
			if err != nil {
				fmt.Print(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		case pathUsers:
			response, err := os.ReadFile("./test-data/users-response.json")
			if err != nil {
				fmt.Print(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(response)

		case pathClients:
			response, err := os.ReadFile("./test-data/client-response.json")
			if err != nil {
				fmt.Print(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		default:
			t.Errorf("Unexpected request path on mock server: %s", r.URL.Path)
		}

	}))
	defer server.Close()

	controller := New(server.URL)
	err := controller.GetControllerInfo()
	if err != nil {
		t.Fatalf("test failure: get controller info")
	}
	err = controller.Login("user", "pass")
	if err != nil {
		t.Fatalf("test failure: controller login: %v", controller)
	}

	err = controller.SetSite("Home")
	if err != nil {
		t.Fatalf("test failure: set site: %v", err)
	}

	clients, err := controller.GetClients()
	if err != nil {
		t.Fatalf("test failure: get clients")
	}

	// assert
	expectedClients := 3
	assert.Len(t, clients, expectedClients)
	assert.Equal(t, "coredns 123", clients[0].Name)
	assert.Equal(t, "10.0.0.5", clients[0].Ip)
	assert.Equal(t, "coredns-123", clients[0].DnsName)

}
