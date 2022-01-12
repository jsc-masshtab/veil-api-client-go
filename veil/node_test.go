package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NodeListGet(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.Node.List()
	assert.Nil(t, err)
	for _, v := range response.Results {
		node, _, err := client.Node.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, node.Id, "", "Node Id can not be empty")
		break
	}

	return
}
