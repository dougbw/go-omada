package omada

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSites(t *testing.T) {

	expectedSites := 2
	assert.Len(t, testController.Sites, expectedSites)
	assert.Equal(t, testData["siteId"], testController.siteId)

}
