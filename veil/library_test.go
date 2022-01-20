package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LibraryListGet(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.Library.List()
	assert.Nil(t, err)
	for _, v := range response.Results {
		entity, _, err := client.Library.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, entity.Id, "", "Library Id can not be empty")

		break
	}

	return
}

func Test_LibraryUpload(t *testing.T) {
	// TODO fix upload files
	t.SkipNow()
	client := NewClient("", "", false)
	response, _, err := client.DataPool.List()
	assert.Nil(t, err)
	if len(response.Results) == 0 {
		t.SkipNow()
	}
	firstDp := response.Results[0]
	entity, _, err := client.Library.Create(firstDp.Id, "base_domain.xml")
	assert.Nil(t, err)
	assert.NotEqual(t, entity.Id, "", "file Id can not be empty")

	return
}

func Test_LibraryDownload(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.Library.List()
	assert.Nil(t, err)
	if len(response.Results) == 0 {
		t.SkipNow()
	}
	for _, v := range response.Results {
		// Больше 50Мб пропускаем
		if v.Size >= 50*1024*1024 {
			continue
		}
		entity, _, err := client.Library.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, entity.Id, "", "file Id can not be empty")
		entity, _, err = client.Library.Download(entity)
		assert.Nil(t, err)
		break
	}

	return
}
