package veil

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	client := NewClient("", "", false)
	response, _, err := client.DataPool.List()
	require.Nil(t, err)
	if len(response.Results) == 0 {
		t.SkipNow()
	}
	firstDp := response.Results[0]
	entity, err := client.Library.Create(firstDp.Id, "base_domain.xml", 0)
	require.Nil(t, err)
	assert.NotEqual(t, entity.Id, "", "file Id can not be empty")

	status, _, err := client.Library.Remove(entity.Id)
	assert.Nil(t, err)
	assert.True(t, status)

	return
}

func Test_LibraryUploadUrl(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.DataPool.List()
	require.Nil(t, err)
	if len(response.Results) == 0 {
		t.SkipNow()
	}
	firstDp := response.Results[0]
	// TODO Check filename exists on datapool
	library, err := client.Library.Create(firstDp.Id, "http://192.168.10.144/test_helper/test_domain.xml", 0)
	assert.Nil(t, err)
	assert.NotEqual(t, library.Id, "", "Library Id can not be empty")
	assert.Equal(t, library.Status, Status.Active, "Library Status should be Active")

	status, _, err := client.Library.Remove(library.Id)
	assert.Nil(t, err)
	assert.True(t, status)

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
