package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UserListGet(t *testing.T) {
	client := NewClient("", "", false)

	response, _, err := client.User.List()
	assert.Nil(t, err)
	for _, v := range response.Results {
		entity, _, err := client.User.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, entity.Id, 0, "User Id can not be empty")
		break
	}

	return
}
