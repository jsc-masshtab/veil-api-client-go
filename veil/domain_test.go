package veil

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

var TestDomainName = NameGenerator("domain")
var TestDomainID = uuid.NewString()
var TestPackerTemplateName = "debian_cloud"
var TestPackerVMName = "packer_vm"
var TestPackerFilename = "debian-11-generic-amd64.qcow2"
var TestPackerIsoFilename = "ubuntu-18.04.3-desktop-amd64.iso"
var TestPackerVdiskName = "packer_vdisk"
var TestPackerVNet = "vms"
var UserDataTemplate = `#cloud-config

groups:
- cloud-users

package_update: false
package_upgrade: false

users:
  - default
  - name: %s
    groups: sudo
    sudo: ALL = (ALL) NOPASSWD:ALL
    shell: /bin/bash
    lock_passwd: false
    plain_text_passwd: %s
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
var baseUser = "user"
var userData = fmt.Sprintf(UserDataTemplate, baseUser, baseUser)

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

func Test_DomainPackerClone(t *testing.T) {
	/*
		Вариант 1. Имитация packer clone
		Условия:
		1. Есть готовый шаблон ВМ с образом cloud_init
		Шаги:
		1. Клонируем шаблон с переводом в ВМ c конфигом cloud_init и включением
		2. Ждем гостевого агента
		3. Подключаемся по ssh
		4. Инициализируем её (provisioning)
		5. Выключаем ВМ
		6. Переводим её в шаблон
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
	cloneConfig := new(DomainCloneConfig)
	cloneConfig.VerboseName = TestPackerVMName
	cloneConfig.Template = false
	cloneConfig.StartOn = true
	cloneConfig.CloudInit = true

	cloudInitConf := new(CloudConfig)
	cloudInitConf.UserData = base64.StdEncoding.EncodeToString([]byte(userData))

	//	metaData := `instance-id: %s
	//local-hostname: %s
	//`
	//metaData = fmt.Sprintf(metaData, newDomain.Id, newDomain.VerboseName)
	//cloudInitConf.MetaData = base64.StdEncoding.EncodeToString([]byte(metaData))
	cloneConfig.CloudInitConfig = *cloudInitConf

	baseTemplate := templatesResponse.Results[0]
	newDomain, _, err := client.Domain.Clone(baseTemplate.Id, *cloneConfig)
	require.Nil(t, err, err)
	assert.Equal(t, newDomain.VerboseName, cloneConfig.VerboseName, fmt.Sprintf("Domain VerboseName should be %s", TestPackerVMName))

	// 2. Ждем гостевого агента
	newDomain, err = newDomain.WaitForGA(client, 80)
	assert.NotEmpty(t, newDomain.GuestUtils.Ipv4, "Domain %s ipV4 in GuestUtils should not be empty", newDomain.VerboseName)

	// 3. Подключаемся по ssh
	sshAddress := newDomain.GuestUtils.Ipv4[0]
	key, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	require.Nil(t, err)
	signer, err := ssh.ParsePrivateKey(key) // Создания подписанта приватного ключа
	require.Nil(t, err)
	sshConfig := &ssh.ClientConfig{
		User: baseUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connection, err := ssh.Dial("tcp", sshAddress+":22", sshConfig)
	defer connection.Close()
	require.Nil(t, err, "Failed to dial: %s", err)
	session, err := connection.NewSession()
	defer session.Close()
	require.Nil(t, err, "Failed to create session: %s", err)

	// 4. Инициализируем её (provisioning)

	cmd := "ls -la /"
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(cmd)
	assert.Nil(t, err)
	fmt.Println(stdoutBuf.String())

	// 5. Выключаем ВМ
	newDomain, _, err = client.Domain.Shutdown(newDomain, true)
	assert.Nil(t, err)

	// 6. Переводим её в шаблон
	newDomain, _, err = client.Domain.Template(newDomain, true)
	assert.Nil(t, err)
	assert.True(t, newDomain.Template)

	// Удаляем ВМ
	status, _, err := client.Domain.Remove(newDomain.Id, true, false)
	assert.Nil(t, err)
	assert.True(t, status)

	return
}

func Test_DomainPackerQcow2(t *testing.T) {
	/*
		Вариант 2. Имитация packer qcow2
		Условия:
		1a. Есть локальный qcow2 с cloud_init на машине, где запускается packer
		1б. Есть диск/файл qcow2 с cloud_init уже в VeiL
		Шаги:
		1. Загружаем qcow2 или выбираем из имеющихся и импортируем его в диск, а также увеличиваем размер
		2. Создаем ВМ с qcow2, конфигом cloud_init и сразу включаем её
		3. Ждем гостевого агента
		4. Подключаемся по ssh
		5. Инициализируем её (provisioning)
		6. Выключаем ВМ
		7. Переводим её в шаблон
	*/
	client := NewClient("", "", false)
	// 1. Загружаем qcow2 или выбираем из имеющихся и импортируем его в диск, а также увеличиваем размер
	filesResponse, _, err := client.Library.ListParams(map[string]string{
		"status":   "ACTIVE",
		"filename": TestPackerFilename,
	})
	require.Nil(t, err, err)
	if len(filesResponse.Results) == 0 {
		t.SkipNow()
	}
	baseFile := filesResponse.Results[0]
	importConfig := new(FileImportConfig)
	importConfig.VerboseName = TestPackerVdiskName
	newVdisk, _, err := client.Library.Import(baseFile.Id, *importConfig)
	assert.NotEqual(t, newVdisk.Id, "", "Vdisk Id can not be empty")
	assert.Equal(t, newVdisk.VerboseName, TestPackerVdiskName, fmt.Sprintf("Vdisk VerboseName should be %s", TestPackerVdiskName))

	_, _, err = client.Vdisk.Extend(newVdisk.Id, 18)
	assert.Nil(t, err)
	newVdisk, err = newVdisk.Refresh(client)
	assert.Nil(t, err)

	// 2. Создаем ВМ с qcow2, конфигом cloud_init и сразу включаем её
	config := new(DomainMultiCreateConfig)

	nodesResponse, _, err := client.Node.List()
	require.Nil(t, err, err)
	if len(nodesResponse.Results) == 0 {
		t.SkipNow()
	}
	randomNode := nodesResponse.Results[rand.Intn(len(nodesResponse.Results))]
	config.Node = randomNode.Id

	vnetsResponse, _, err := client.Vnet.ListParams(map[string]string{
		"status":       "ACTIVE",
		"node":         randomNode.Id,
		"verbose_name": TestPackerVNet,
	})
	require.Nil(t, err, err)
	if len(vnetsResponse.Results) == 0 {
		t.SkipNow()
	}
	vnet := vnetsResponse.Results[0]
	config.VmachineInfs = []VMachineInfSoftCreate{
		{
			Vnetwork:  vnet.Id,
			NicDriver: "virtio",
		},
	}

	config.Vdisks = []VdiskAttach{
		{
			Vdisk: newVdisk.Id,
		},
	}

	config.DomainId = TestDomainID
	config.VerboseName = TestPackerVMName
	config.MemoryCount = 4096
	config.StartOn = true

	config.CloudInit = true
	cloudInitConf := new(CloudConfig)
	cloudInitConf.UserData = base64.StdEncoding.EncodeToString([]byte(userData))
	config.CloudInitConfig = *cloudInitConf

	newDomain, _, err := client.Domain.MultiCreate(*config)
	require.Nil(t, err, err)
	assert.NotEqual(t, newDomain.Id, "", "Domain Id can not be empty")

	// 3. Ждем гостевого агента
	newDomain, err = newDomain.WaitForGA(client, 80)
	assert.NotEmpty(t, newDomain.GuestUtils.Ipv4, "Domain %s ipV4 in GuestUtils should not be empty", newDomain.VerboseName)

	// 4. Подключаемся по ssh
	sshAddress := newDomain.GuestUtils.Ipv4[0]
	key, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	require.Nil(t, err)
	signer, err := ssh.ParsePrivateKey(key) // Создания подписанта приватного ключа
	require.Nil(t, err)
	sshConfig := &ssh.ClientConfig{
		User: baseUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connection, err := ssh.Dial("tcp", sshAddress+":22", sshConfig)
	defer connection.Close()
	require.Nil(t, err, "Failed to dial: %s", err)
	session, err := connection.NewSession()
	defer session.Close()
	require.Nil(t, err, "Failed to create session: %s", err)

	// 5. Инициализируем её (provisioning)

	cmd := "ls -la /"
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(cmd)
	assert.Nil(t, err)
	fmt.Println(stdoutBuf.String())

	// 6. Выключаем ВМ
	newDomain, _, err = client.Domain.Shutdown(newDomain, true)
	assert.Nil(t, err)

	// 6. Переводим её в шаблон
	newDomain, _, err = client.Domain.Template(newDomain, true)
	assert.Nil(t, err)
	assert.True(t, newDomain.Template)

	// Удаляем ВМ
	status, _, err := client.Domain.Remove(newDomain.Id, true, false)
	assert.Nil(t, err)
	assert.True(t, status)
	return
}

func Test_DomainPackerIso(t *testing.T) {
	/*
		Вариант 2. Имитация packer iso
		Условия:
		1a. Есть локальный образ с cloud_init на машине, где запускается packer
		1б. Есть образ с cloud_init для загрузки по url
		1в. Есть образ с cloud_init уже в VeiL
		1г. Есть диск/файл qcow2 с cloud_init уже в VeiL
		Шаги:
		1. Загружаем образ/qcow2 или выбираем из имеющихся
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
