package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VMachineInfListGet(t *testing.T) {
	client := NewClient("", "", false)

	response, _, err := client.VMachineInf.List()
	assert.Nil(t, err)
	for _, v := range response.Results {
		entity, _, err := client.VMachineInf.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, entity.Id, "", "VMachineInf Id can not be empty")

		entity, err = entity.Refresh(client)
		assert.Nil(t, err)
		break
	}

	return
}
