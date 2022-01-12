package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Task(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.Task.List()
	assert.Nil(t, err)
	for _, v := range response.Results {
		task, _, err := client.Task.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, task.Id, "", "Task Id can not be empty")

		break
	}

	return

}
