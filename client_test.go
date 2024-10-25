package omada

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClients(t *testing.T) {

	clients, err := testController.GetClients()
	if err != nil {
		t.Fatalf("test failure on 'GetClients': %v", err)
	}

	expectedClients := 5
	assert.Equal(t, testController.siteId, "Default")
	assert.Len(t, clients, expectedClients)

	assert.Equal(t, "Client 001", clients[1].Name)
	assert.Equal(t, "client-001", clients[1].DnsName)

	assert.Equal(t, "Google Nest Mini", clients[2].Name)
	assert.Equal(t, "google-nest-mini", clients[2].DnsName)

	// if client name matches mac address then use hostname as DNS name
	assert.Equal(t, "AA-AA-AA-AA-AA-05", clients[3].Name)
	assert.Equal(t, "WIN-XABHS9AHQZZ", clients[3].HostName)
	assert.Equal(t, "win-xabhs9ahqzz", clients[3].DnsName)

	// if client name matches mac address but hostname is not set then use client name
	assert.Equal(t, "AA-AA-AA-AA-AA-04", clients[0].Name)
	assert.Equal(t, "--", clients[0].HostName)
	assert.Equal(t, "aa-aa-aa-aa-aa-04", clients[0].DnsName)

	ip := net.ParseIP(clients[0].Ip) // ParseIP returns nil rather than error if unable to parse
	assert.NotNil(t, ip)
}
