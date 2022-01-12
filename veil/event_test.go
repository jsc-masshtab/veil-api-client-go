package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Event(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.Event.List()
	assert.Nil(t, err)
	for _, v := range response.Results {
		task, _, err := client.Event.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, task.Id, "", "Event Id can not be empty")

		break
	}

	return

}
