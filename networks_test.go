package omada

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworks(t *testing.T) {

	networks, err := testController.GetNetworks()
	if err != nil {
		t.Fatalf("test failure on 'GetNetworks': %v", err)
	}

	for _, network := range networks {
		slog.Debug("Id:         " + network.Id)
		slog.Debug("Name:       " + network.Name)
		slog.Debug("Domain:     " + network.Domain)
		slog.Debug("Purpose:    " + network.Purpose)
		slog.Debug("Subnet:     " + network.Subnet)
		slog.Debug("Ipv6Subnet: " + network.Ipv6Subnet)

	}

	expectedCount := 4
	assert.Len(t, networks, expectedCount)

	assert.Equal(t, "5efee9f40e4a4b066b212f89", networks[0].Id)
	assert.Equal(t, "Default", networks[0].Name)
	assert.Equal(t, "omada.home", networks[0].Domain)
	assert.Equal(t, "interface", networks[0].Purpose)
	assert.Equal(t, "10.0.0.1/24", networks[0].Subnet)
	assert.Equal(t, "2001:db8:f7::/64", networks[0].Ipv6Subnet)

	assert.Equal(t, "65bee8cd86967408e454daac", networks[1].Id)
	assert.Equal(t, "Guest", networks[1].Name)
	assert.Equal(t, "", networks[1].Domain)
	assert.Equal(t, "interface", networks[1].Purpose)
	assert.Equal(t, "10.0.200.1/24", networks[1].Subnet)
	assert.Equal(t, "", networks[1].Ipv6Subnet)

	assert.Equal(t, "66b13848e2f8c96db92bca15", networks[2].Id)
	assert.Equal(t, "Test", networks[2].Name)
	assert.Equal(t, "", networks[2].Domain)
	assert.Equal(t, "vlan", networks[2].Purpose)
	assert.Equal(t, "", networks[2].Subnet)
	assert.Equal(t, "", networks[2].Ipv6Subnet)

	assert.Equal(t, "62238f02dd3c1a436c3e0a64", networks[3].Id)
	assert.Equal(t, "Work", networks[3].Name)
	assert.Equal(t, "omada.work", networks[3].Domain)
	assert.Equal(t, "interface", networks[3].Purpose)
	assert.Equal(t, "10.0.100.1/24", networks[3].Subnet)
	assert.Equal(t, "", networks[3].Ipv6Subnet)

}
