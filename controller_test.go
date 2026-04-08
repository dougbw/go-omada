package omada

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testController Controller
var testData = map[string]string{}

func TestMain(m *testing.M) {

	testData["controllerId"] = "123bee230c77bbb45d9c8545d04d700a"
	testData["siteId"] = "Default"

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

func serveFile(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func newTestServer(routes map[string]http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, ok := routes[r.URL.Path]
		if !ok {
			log.Fatalf("unexpected request path on mock server: %s", r.URL.Path)
		}
		handler(w, r)
	}))
}

func commonRoutes(infoFile string) map[string]http.HandlerFunc {
	controllerId := testData["controllerId"]
	return map[string]http.HandlerFunc{
		"/api/info": serveFile(infoFile),
		fmt.Sprintf("/%s/api/v2/login", controllerId):         serveFile("./test-data/login-response.json"),
		fmt.Sprintf("/%s/api/v2/users/current", controllerId): serveFile("./test-data/users-response.json"),
	}
}

func setupTestServer() *httptest.Server {
	controllerId := testData["controllerId"]
	siteId := testData["siteId"]

	routes := commonRoutes("./test-data/info-response.json")
	routes[fmt.Sprintf("/%s/api/v2/sites/%s/clients", controllerId, siteId)] = serveFile("./test-data/clients-response.json")
	routes[fmt.Sprintf("/%s/api/v2/sites/%s/devices", controllerId, siteId)] = serveFile("./test-data/devices-response.json")
	routes[fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/networks", controllerId, siteId)] = serveFile("./test-data/networks-response.json")
	routes[fmt.Sprintf("/%s/api/v2/sites/%s/setting/service/dhcp", controllerId, siteId)] = serveFile("./test-data/dhcp-reservation-response.json")

	return newTestServer(routes)
}

func TestLogin(t *testing.T) {
	assert.Equal(t, testData["controllerId"], testController.controllerId)
}

func setupOpenAPITestServer(t *testing.T) (*httptest.Server, Controller) {
	t.Helper()
	controllerId := testData["controllerId"]
	siteId := testData["siteId"]

	routes := commonRoutes("./test-data/openapi/info-response.json")
	routes[fmt.Sprintf("/openapi/v2/%s/sites/%s/clients", controllerId, siteId)] = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "web-local", r.Header.Get("Omada-Request-Source"))
		assert.NotEmpty(t, r.Header.Get("Csrf-Token"))
		serveFile("./test-data/clients-response.json")(w, r)
	}

	server := newTestServer(routes)
	c := New(server.URL)
	if err := c.GetControllerInfo(); err != nil {
		t.Fatal("GetControllerInfo:", err)
	}
	if err := c.Login("user", "pass"); err != nil {
		t.Fatal("Login:", err)
	}
	if err := c.SetSite("Home"); err != nil {
		t.Fatal("SetSite:", err)
	}
	return server, c
}
