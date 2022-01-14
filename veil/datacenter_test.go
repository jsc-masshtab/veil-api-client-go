package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DatacenterListGet(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.DataCenter.List()
	assert.Nil(t, err)
	for _, v := range response.Results {
		node, _, err := client.DataCenter.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, node.Id, "", "DataCenter Id can not be empty")
		break
	}

	return
}
