package veil

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

var TestDomainName = NameGenerator("domain")
var TestDomainID = uuid.NewString()
var TestPackerTemplateName = "debian_cloud"

func Test_DomainList(t *testing.T) {
	client := NewClient("", "", false)

	_, _, err := client.Domain.List()
	assert.Nil(t, err)

	return
}

func Test_DomainCreate(t *testing.T) {
	client := NewClient("", "", false)
	config := new(DomainCreateConfig)
	config.DomainId = TestDomainID
	config.VerboseName = TestDomainName
	config.MemoryCount = 50
	domain, _, err := client.Domain.Create(*config)
	require.Nil(t, err)
	assert.NotEqual(t, domain.Id, "", "Domain Id can not be empty")

	return
}

func Test_DomainGet(t *testing.T) {
	client := NewClient("", "", false)

	domain, _, err := client.Domain.Get(TestDomainID)
	assert.Nil(t, err)
	assert.NotEqual(t, domain.Id, "", "Domain Id can not be empty")

	return
}

func Test_DomainPower(t *testing.T) {
	client := NewClient("", "", false)

	domain, _, err := client.Domain.Get(TestDomainID)
	assert.Nil(t, err)
	assert.NotEqual(t, domain.Id, "", "Domain Id can not be empty")
	domain, _, err = client.Domain.Start(domain)
	assert.Nil(t, err)
	domain, _, err = client.Domain.Suspend(domain)
	assert.Nil(t, err)
	domain, _, err = client.Domain.Resume(domain)
	assert.Nil(t, err)
	domain, _, err = client.Domain.Reboot(domain, true)
	assert.Nil(t, err)
	domain, _, err = client.Domain.Shutdown(domain, true)
	assert.Nil(t, err)
	domain, _, err = client.Domain.Template(domain, true)
	assert.Nil(t, err)
	assert.True(t, domain.Template)
	domain, _, err = client.Domain.Template(domain, false)
	assert.Nil(t, err)
	assert.False(t, domain.Template)

	return
}

func Test_DomainUpdate(t *testing.T) {
	client := NewClient("", "", false)
	config := new(DomainUpdateConfig)
	newName := "test"
	config.VerboseName = newName
	domain, _, err := client.Domain.Update(TestDomainID, *config)
	domain.Refresh(client)
	assert.Nil(t, err)
	assert.NotEqual(t, domain.VerboseName, newName, "Domain VerboseName should be test")

	return
}

func Test_DomainRemove(t *testing.T) {
	client := NewClient("", "", false)

	status, _, err := client.Domain.Remove(TestDomainID, true, false)
	assert.Nil(t, err)
	assert.True(t, status)

	return
}

func Test_DomainMultiCreate(t *testing.T) {
	client := NewClient("", "", false)

	nodesResponse, _, err := client.Node.List()
	require.Nil(t, err, err)
	if len(nodesResponse.Results) == 0 {
		t.SkipNow()
	}
	randomNode := nodesResponse.Results[rand.Intn(len(nodesResponse.Results))]

	config := new(DomainMultiCreateConfig)
	config.DomainId = TestDomainID
	config.VerboseName = TestDomainName
	config.Node = randomNode.Id
	domain, _, err := client.Domain.MultiCreate(*config)
	require.Nil(t, err, err)
	assert.NotEqual(t, domain.Id, "", "Domain Id can not be empty")
	status, _, err := client.Domain.Remove(TestDomainID, true, false)
	assert.Nil(t, err)
	assert.True(t, status)

	return
}

func Test_DomainMultiCreateThin(t *testing.T) {
	client := NewClient("", "", false)

	templatesResponse, _, err := client.Domain.ListParams(map[string]string{
		"template": "true",
		"status":   "ACTIVE",
	})
	require.Nil(t, err, err)
	if len(templatesResponse.Results) == 0 {
		t.SkipNow()
	}
	randomTemplate := templatesResponse.Results[rand.Intn(len(templatesResponse.Results))]

	config := new(DomainMultiCreateConfig)
	config.DomainId = TestDomainID
	config.VerboseName = TestDomainName
	config.Parent = randomTemplate.Id
	config.Thin = true
	config.StartOn = true
	domain, _, err := client.Domain.MultiCreate(*config)
	require.Nil(t, err, err)
	assert.NotEqual(t, domain.Id, "", "Domain Id can not be empty")
	status, _, err := client.Domain.Remove(TestDomainID, true, false)
	assert.Nil(t, err)
	assert.True(t, status)

	return
}

func Test_DomainPacker_1(t *testing.T) {
	/*
		Вариант 1. Имитация packer clone
		Условия:
		1. Есть готовый шаблон ВМ с образом cloud_init
		Шаги:
		1. Клонируем шаблон с переводом в ВМ
		2. Добавляем к ней конфиг cloud_init
		3. Включаем ВМ
		4. Ждем гостевого агента
		5. Подключаемся по ssh
		6. Инициализируем её (provisioning)
		7. Выключаем ВМ
		8. Переводим её в шаблон
	*/
	client := NewClient("", "", false)
	// check templates
	templatesResponse, _, err := client.Domain.ListParams(map[string]string{
		"template": "true",
		"status":   "ACTIVE",
		"name":     TestPackerTemplateName,
	})
	require.Nil(t, err, err)
	if len(templatesResponse.Results) == 0 {
		t.SkipNow()
	}
	// 1. Клонируем шаблон с переводом в ВМ
	domainName := "packer_vm"
	cloneConfig := new(DomainCloneConfig)
	cloneConfig.VerboseName = domainName
	cloneConfig.Template = false
	baseTemplate := templatesResponse.Results[0]
	newDomain, _, err := client.Domain.Clone(baseTemplate.Id, *cloneConfig)
	require.Nil(t, err, err)
	assert.Equal(t, newDomain.VerboseName, cloneConfig.VerboseName, fmt.Sprintf("Domain VerboseName should be %s", domainName))
	// 2. Добавляем к ней конфиг cloud_init
	cloudInitConf := new(CloudConfig)
	userData := `#cloud-config

groups:
- cloud-users

package_update: false
package_upgrade: false

users:
  - default
  - name: user
    groups: sudo
    sudo: ALL = (ALL) NOPASSWD:ALL
    shell: /bin/bash
    lock_passwd: false
    plain_text_passwd: user
    # passwd: $6$rounds = 4096$9cYh.jYsend9bOZ$VBqFtH6Jc6cgpYga.sWD.G5l/h.Fedn.CRO7ouw7S7JiMbwXvf5cuENpOk9W4pqAAmF7vxKJy62QCHZ9xVvAd0
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCWRyWsZAFDAiUuCPAz0cN1jFREnnBdLkpIoQygLJqkrzh85l47Omib+IlryIaa7QsjaKhI2dYUfviTiOeUM0yLH17YD7IR+8n2uADy3kNHbjwDn3+9OOfCGExLXH6Az1imenWJj6ErLmelTJi66xLWcGQhBNtr37XwOlL8eguP4TwZ1LmoUqWseKXEerUoOKqP2abYu5zgWNtkWJ5604V8lvQt5JgMJMqr7oGCIT/DgD/ndqOOu0G6698deEk/ooADVB1CUglrPni+ZPBHhwwMrovpkKgwbOTUXrmE5I9OrmsjLGiaLkjsSyQrfx5xfrXhogbCE174PWJaCy8zD7HLGArmhsBnMz8FKEbX/We547llCKGPGmc4H6IMhbryiZky3XuGK3nBKvmOiwwKUNoamt7yXUIRfFcoOqhC63DZfHT/4OvfvKnv3HtnY2VoDgZCaYCcT6ZZwntk2p6LY2zoDwqThXLtvZwouPkhtdOs2ATvW04CMnXCKBsu2W76c60= user@user-To-be-filled-by-O-E-M

manage_resolv_conf: false
manage_etc_hosts: false

packages:
  - qemu-guest-agent
  - openssl
  - curl
  - openssh-server

runcmd:
  - systemctl start qemu-guest-agent

final_message: "The system is finally up, after $UPTIME seconds"`
	sDec := base64.StdEncoding.EncodeToString([]byte(userData))
	fmt.Println(sDec)
	cloudInitConf.MetaData = "Cmluc3RhbmNlLWlkOiBiNzNlODZkNi02Yjc3LTRkMTQtYmNlZi1hMWRmNzM2Y2U4N2UKbG9jYWwtaG9zdG5hbWU6IHNvbWVkb21haW4KCg=="
	cloudInitConf.UserData = sDec
	cloudInit := new(CloudInitConfig)
	cloudInit.CloudInit = true
	cloudInit.CloudInitConfig = *cloudInitConf
	newDomain, _, err = client.Domain.CloudInit(newDomain, *cloudInit)
	assert.Nil(t, err)
	// 3. Включаем ВМ
	_, _, err = client.Domain.Start(newDomain)
	assert.Nil(t, err)
	// 4. Ждем гостевого агента
	newDomain, err = newDomain.WaitForGA(client, 80)
	assert.Nil(t, err)

	// Удаляем ВМ
	status, _, err := client.Domain.Remove(newDomain.Id, true, false)
	assert.Nil(t, err)
	assert.True(t, status)

	return
}

func Test_DomainPacker_2(t *testing.T) {
	/*
		Вариант 2. Имитация packer iso
		Условия:
		1a. Есть локальный образ с cloud_init на машине, где запускается packer
		1б. Есть образ с cloud_init для загрузки по url
		1в. Есть образ с cloud_init уже в VeiL
		Шаги:
		1. Загружаем образ или выбираем из имеющихся
		2. Создаем ВМ с образом cloud_init
		3. Включаем ВМ
		4. Ждем гостевого агента
		5. Подключаемся по ssh
		6. Инициализируем её (provisioning)
		7. Выключаем ВМ
		8. Переводим её в шаблон

	*/
	t.SkipNow()
	//client := NewClient("", "", false)

	return
}
