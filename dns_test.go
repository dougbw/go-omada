package omada

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDns(t *testing.T) {

	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "Test 01", want: "test-01"},
		{input: "Test-01", want: "test-01"},
		{input: "Test_01", want: "test01"},
		{input: "Test.01", want: "test01"},
	}

	for _, v := range tests {
		output := makeDNSSafe(v.input)
		assert.Equal(t, v.want, output)
	}

}
