package omada

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDhcp(t *testing.T) {

	reservations, err := testController.GetDhcpReservations()
	if err != nil {
		t.Fatalf("test failure on 'GetDhcpReservations': %v", err)
	}

	expectedCount := 3
	assert.Len(t, reservations, expectedCount)
	assert.Equal(t, "AA-BB-CC3-0C-8B-C6", reservations[0].ClientName)
	assert.Equal(t, "10.0.0.100", reservations[0].IP)

}
