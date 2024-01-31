package vmcreator

import (
	"fmt"

	"github.com/tigerinus/libvirt-go-demo/config"
	"github.com/tigerinus/libvirt-go-demo/installermedia"
	"github.com/tigerinus/libvirt-go-demo/util"
	"github.com/tigerinus/libvirt-go-demo/vmconfigurator"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type VMCreator struct {
	InstallMedia *installermedia.InstallerMedia
	connection   *libvirt.Connect
}

func NewVMCreator(installMedia *installermedia.InstallerMedia) *VMCreator {
	return &VMCreator{
		InstallMedia: installMedia,
	}
}

func (vmc *VMCreator) createDomainConfig(name, title, volumePath string) (*libvirtxml.Domain, error) {
	capsXML, err := vmc.connection.GetCapabilities()
	if err != nil {
		return nil, err
	}

	var caps libvirtxml.Caps
	if err2 := caps.Unmarshal(capsXML); err2 != nil {
		return nil, err2
	}

	domcapsXML, err := vmc.connection.GetDomainCapabilities("", "", "", "", 0)
	if err != nil {
		return nil, err
	}

	var domcaps libvirtxml.DomainCaps
	if err3 := domcaps.Unmarshal(domcapsXML); err3 != nil {
		return nil, err3
	}

	config, err := vmconfigurator.CreateDomainConfig(vmc.InstallMedia, volumePath, caps, domcaps)
	if err != nil {
		return nil, err
	}

	config.Name = name
	config.Title = title

	return config, nil
}

func (vmc *VMCreator) CreateVM(clone bool) (*libvirt.Domain, error) {
	name, title, err := vmc.createDomainNameAndTitleFromMedia()
	if err != nil {
		return nil, err
	}

	vmc.InstallMedia.PrepareForInstallation(name)

	var volumePath *string
	if vmc.InstallMedia.SkipImport && !clone {
		volumePath = &vmc.InstallMedia.DeviceFile

		fmt.Printf("Skipping import. Using '%s' as target volume\n", *volumePath)
	} else {
		volume, err2 := vmc.createTargetVolume(name, config.DefaultStorage)
		if err2 != nil {
			return nil, err2
		}
		_volumePath, err2 := volume.GetPath()
		if err2 != nil {
			return nil, err2
		}

		volumePath = &_volumePath
	}

	config, err := vmc.createDomainConfig(name, title, *volumePath)
	if err != nil {
		return nil, err
	}

	xmlConfig, err := config.Marshal()
	if err != nil {
		return nil, err
	}

	return vmc.connection.DomainDefineXML(xmlConfig)
}

func (vmc *VMCreator) createDomainNameAndTitleFromMedia() (name string, title string, err error) {
	baseName, baseTitle := vmc.createDomainBaseNameAndTitle()

	name = baseName
	title = baseTitle

	pool, err := EnsureStoragePool(vmc.connection)

	for i := 2; ; i++ {
		domain, err := vmc.connection.LookupDomainByName(name)
		if err != nil {
			return "", "", err
		}

		if domain == nil {
			break
		}

		domain, err = vmc.connection.LookupDomainByName(title)
		if err != nil {
			return "", "", err
		}

		if domain == nil {
			break
		}

		volume, err := pool.LookupStorageVolByName(name)
		if err != nil {
			return "", "", err
		}

		if volume == nil {
			break
		}

		name = fmt.Sprintf("%s-%d", baseName, i)
		title = fmt.Sprintf("%s-%d", baseTitle, i)
	}

	return
}

func (vmc *VMCreator) createDomainBaseNameAndTitle() (baseName string, baseTitle string) {
	baseTitle = vmc.InstallMedia.Label
	baseName = "unknown"
	return
}

func (vmc *VMCreator) createTargetVolume(name string, storage uint64) (*libvirt.StorageVol, error) {
	pool, err := EnsureStoragePool(vmc.connection)
	if err != nil {
		return nil, err
	}

	config := vmconfigurator.CreateVolumeConfig(name, storage)
	xmlConfig, err := config.Marshal()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Creating volume '%s' ..\n", name)

	volume, err := pool.StorageVolCreateXML(xmlConfig, 0)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created volume '%s'.\n", name)

	return volume, nil
}

func EnsureStoragePool(connection *libvirt.Connect) (*libvirt.StoragePool, error) {
	pool, err := GetStoragePool(*connection)
	if err != nil {
		return nil, err
	}

	if pool == nil {
		fmt.Println("Creating storage pool..")
		poolConfig := vmconfigurator.GetPoolConfig()
		xmlConfig, err3 := poolConfig.Marshal()
		if err3 != nil {
			return nil, err3
		}

		pool, err = connection.StoragePoolCreateXML(xmlConfig, libvirt.STORAGE_POOL_CREATE_NORMAL)
		if err != nil {
			return nil, err
		}

		if err2 := pool.Build(libvirt.STORAGE_POOL_BUILD_NEW); err2 != nil {
			return nil, err2
		}

		fmt.Println("Created storage pool.")
	}

	poolPath := util.GetUserPkgData("images")
	util.EnsureDirectory(poolPath)

	poolInfo, err := pool.GetInfo()
	if err != nil {
		return nil, err
	}

	if poolInfo.State == libvirt.STORAGE_POOL_INACTIVE {
		if err := pool.Create(libvirt.STORAGE_POOL_CREATE_NORMAL); err != nil {
			return nil, err
		}
	}

	if err := pool.Refresh(0); err != nil {
		return nil, err
	}

	if err := pool.SetAutostart(true); err != nil {
		return nil, err
	}

	return pool, nil
}

func GetStoragePool(connection libvirt.Connect) (*libvirt.StoragePool, error) {
	return connection.LookupStoragePoolByName(config.PackageTarname)
}
