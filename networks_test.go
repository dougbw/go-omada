package omada

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworks(t *testing.T) {

	networks, err := testController.GetNetworks()
	if err != nil {
		t.Fatalf("test failure on 'GetNetworks': %v", err)
	}

	expectedCount := 4
	assert.Len(t, networks, expectedCount)
	assert.Equal(t, "Default", networks[0].Name)
	assert.Equal(t, "omada.home", networks[0].Domain)
	assert.Equal(t, "10.0.0.1/24", networks[0].Subnet)
	assert.Equal(t, 1, networks[0].Vlan)
	assert.Equal(t, "Guest", networks[1].Name)
	assert.Equal(t, 200, networks[1].Vlan)

}
