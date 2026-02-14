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

	expectedClients := 3
	assert.Equal(t, testController.siteId, "Default")
	assert.Len(t, clients, expectedClients)

	assert.Equal(t, "Client 001", clients[0].Name)
	assert.Equal(t, "client-001", clients[0].DnsName)
	ip := net.ParseIP(clients[0].Ip) // ParseIP returns nil rather than error if unable to parse
	assert.NotNil(t, ip)

	assert.Equal(t, "Google Nest Mini", clients[1].Name)
	assert.Equal(t, "google-nest-mini", clients[1].DnsName)
	assert.Equal(t, 1, len(clients[1].IPV6List))
	assert.Equal(t, "2001:db8:f7::103", clients[1].IPV6List[0])

	assert.Equal(t, "win10-vm", clients[2].Name)
	assert.Equal(t, "win10-vm", clients[2].DnsName)
	assert.Equal(t, 1, len(clients[2].IPV6List))
	assert.Equal(t, "2001:db8:f7::102", clients[2].IPV6List[0])
}
