package omada

import (
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClients(t *testing.T) {

	clients, err := testController.GetClients()
	if err != nil {
		t.Fatalf("test failure on 'GetClients': %v", err)
	}

	assert.Equal(t, testController.siteId, "Default")
	expectedClients, _ := strconv.Atoi(testData["expectedClients"])
	assert.Len(t, clients, expectedClients)
	assert.Equal(t, "Client 001", clients[0].Name)
	assert.Equal(t, "Google Nest Mini", clients[1].Name)
	assert.Equal(t, "client-001", clients[0].DnsName)
	assert.Equal(t, "google-nest-mini", clients[1].DnsName)

	ip := net.ParseIP(clients[0].Ip) // ParseIP returns nil rather than error if unable to parse
	assert.NotNil(t, ip)
}
