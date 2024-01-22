package omada

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testController Controller
var testData = map[string]string{}

func TestMain(m *testing.M) {

	testData["controllerId"] = "123bee230c77bbb45d9c8545d04d700a"
	testData["siteId"] = "Default"
	testData["expectedClients"] = "3"
	testData["expectedSites"] = "2"

	testServer := setupTestServer()
	defer testServer.Close()

	testController = New(testServer.URL)
	err := testController.GetControllerInfo()
	if err != nil {
		log.Fatalf("test failure on 'GetControllerInfo': %v", err)
	}
	err = testController.Login("user", "pass")
	if err != nil {
		log.Fatalf("test failure on 'Login': %v", err)
	}
	err = testController.SetSite("Home")
	if err != nil {
		log.Fatalf("test failure on 'SetSite': %v", err)
	}
	run := m.Run()
	os.Exit(run)

}

func setupTestServer() *httptest.Server {

	controllerId := testData["controllerId"]
	siteId := testData["siteId"]
	pathLogin := fmt.Sprintf("/%s/api/v2/login", controllerId)
	pathUsers := fmt.Sprintf("/%s/api/v2/users/current", controllerId)
	pathClients := fmt.Sprintf("/%s/api/v2/sites/%s/clients", controllerId, siteId)
	pathDevices := fmt.Sprintf("/%s/api/v2/sites/%s/devices", controllerId, siteId)
	pathNetworks := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/networks", controllerId, siteId)

	responses := map[string]string{
		"/api/info":  "./test-data/info-response.json",
		pathLogin:    "./test-data/login-response.json",
		pathUsers:    "./test-data/users-response.json",
		pathClients:  "./test-data/clients-response.json",
		pathDevices:  "./test-data/devices-response.json",
		pathNetworks: "./test-data/networks-response.json",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseFile, ok := responses[r.URL.Path]
		if !ok {
			log.Fatalf("Unexpected request path on mock server: %s", r.URL.Path)
		}
		response, err := os.ReadFile(responseFile)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}))

	return server
}

func TestLogin(t *testing.T) {

	testServer := setupTestServer()
	defer testServer.Close()

	controller := New(testServer.URL)
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

	expectedSites, _ := strconv.Atoi(testData["expectedSites"])
	assert.Len(t, controller.Sites, expectedSites)
	assert.Equal(t, testData["controllerId"], testController.controllerId)
	assert.Equal(t, testData["siteId"], testController.siteId)

}
