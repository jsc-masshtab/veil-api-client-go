package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsoListGet(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.Iso.List()
	assert.Nil(t, err)
	for _, v := range response.Results {
		entity, _, err := client.Iso.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, entity.Id, "", "Iso Id can not be empty")

		break
	}

	return

}

func Test_IsoUpload(t *testing.T) {
	t.SkipNow()
	client := NewClient("", "", false)
	response, _, err := client.DataPool.List()
	assert.Nil(t, err)
	if len(response.Results) == 0 {
		t.SkipNow()
	}
	firstDp := response.Results[0]
	iso, _, err := client.Iso.Create(firstDp.Id, "test_live.iso")
	assert.Nil(t, err)
	assert.NotEqual(t, iso.Id, "", "Iso Id can not be empty")

	return

}

func Test_IsoDownload(t *testing.T) {
	client := NewClient("", "", false)
	response, _, err := client.Iso.List()
	assert.Nil(t, err)
	if len(response.Results) == 0 {
		t.SkipNow()
	}
	for _, v := range response.Results {
		// Больше 50Мб пропускаем
		if v.Size >= 50*1024*1024 {
			continue
		}
		iso, _, err := client.Iso.Get(v.Id)
		assert.Nil(t, err)
		assert.NotEqual(t, iso.Id, "", "Iso Id can not be empty")
		iso, _, err = client.Iso.Download(iso)
		assert.Nil(t, err)
		break
	}

	return

}
