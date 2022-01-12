package veil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VdiskCreateSync(t *testing.T) {
	client := NewClient("", "", false)
	dpResponse, _, err := client.DataPool.List()
	assert.Nil(t, err)
	if len(dpResponse.Results) == 0 {
		t.SkipNow()
	}
	firstDp := dpResponse.Results[0]
	config := new(VdiskCreate)
	config.VerboseName = NameGenerator("vdisk")
	config.Preallocation = false
	config.Datapool = firstDp.Id
	config.Size = 0.1
	vdisk, _, err := client.Vdisk.Create(config, false)
	assert.Nil(t, err)
	vdisk, _, err = client.Vdisk.Get(vdisk.Id)
	assert.Nil(t, err)

	vdiskResponseBase, _, err := client.Vdisk.List()
	assert.Nil(t, err)
	if len(vdiskResponseBase.Results) == 0 {
		t.Errorf("no vdisks")
	}

	vdiskResponse, _, err := client.Vdisk.ListParams(map[string]string{})
	assert.Nil(t, err)
	if len(vdiskResponse.Results) == 0 {
		t.Errorf("no vdisks")
	}
	vdisk, _, err = client.Vdisk.Update(vdisk.Id, "test")
	assert.Nil(t, err)

	status, _, err := client.Vdisk.Remove(vdisk.Id)
	assert.Nil(t, err)

	if !status {
		t.Error(err)
		return
	}
	return
}

func Test_VdiskCreateAsync(t *testing.T) {
	client := NewClient("", "", false)
	dpResponse, _, err := client.DataPool.List()
	assert.Nil(t, err)
	if len(dpResponse.Results) == 0 {
		t.SkipNow()
	}
	firstDp := dpResponse.Results[0]
	config := new(VdiskCreate)
	config.VerboseName = NameGenerator("vdisk")
	config.Preallocation = false
	config.Datapool = firstDp.Id
	config.Size = 0.1
	vdisk, _, err := client.Vdisk.Create(config, true)
	assert.Nil(t, err)
	vdisk, _, err = client.Vdisk.Get(vdisk.Id)
	assert.Nil(t, err)

	vdiskResponse, _, err := client.Vdisk.ListParams(map[string]string{})
	assert.Nil(t, err)
	if len(vdiskResponse.Results) == 0 {
		t.Errorf("no vdisks")
	}
	vdisk, _, err = client.Vdisk.Update(vdisk.Id, "test")
	assert.Nil(t, err)

	status, _, err := client.Vdisk.Remove(vdisk.Id)
	assert.Nil(t, err)

	if !status {
		t.Error(err)
		return
	}
	return
}
