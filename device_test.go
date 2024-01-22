package omada

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDevices(t *testing.T) {

	devices, err := testController.GetDevices()
	if err != nil {
		t.Fatalf("test failure on 'GetDevices': %v", err)
	}
	expectedDevices := 6
	assert.Len(t, devices, expectedDevices)
	assert.Equal(t, "access-point-01", devices[0].Name)
	assert.Equal(t, "access-point-02", devices[1].Name)

}
