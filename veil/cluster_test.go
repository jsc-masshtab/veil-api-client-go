package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ClusterListGet(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.Cluster.List()
	assert.Nil(t, err)
	for _, v := range response.Results {
		node, _, err := client.Cluster.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, node.Id, "", "Cluster Id can not be empty")
		break
	}

	return
}
