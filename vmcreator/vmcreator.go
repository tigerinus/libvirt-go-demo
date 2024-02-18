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

	config, err := vmconfigurator.CreateDomainConfig(vmc.InstallMedia, volumePath, caps)
	if err != nil {
		return nil, err
	}

	config.Name = name
	config.Title = title

	return config, nil
}

func (vmc *VMCreator) CreateVM(clone bool) (*libvirt.Domain, error) {
	if vmc.connection == nil {
		conn, err := libvirt.NewConnect("qemu:///session")
		if err != nil {
			return nil, err
		}

		vmc.connection = conn
	}

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
		volume, err2 := vmc.createTargetVolume(name, uint64(config.DefaultStorage))
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

func (vmc *VMCreator) createDomainNameAndTitleFromMedia() (string, string, error) {
	baseName, baseTitle := vmc.createDomainBaseNameAndTitle()

	name := baseName
	title := baseTitle

	pool, err := vmc.EnsureStoragePool()
	if err != nil {
		return "", "", err
	}

	for i := 2; ; i++ {
		if _, err = vmc.connection.LookupDomainByName(name); err != nil {
			if libvirtErr, ok := err.(libvirt.Error); ok && libvirtErr.Code == libvirt.ERR_NO_DOMAIN {
				break
			}

			return "", "", err
		}

		if _, err = vmc.connection.LookupDomainByName(title); err != nil {
			if libvirtErr, ok := err.(libvirt.Error); ok && libvirtErr.Code == libvirt.ERR_NO_DOMAIN {
				break
			}

			return "", "", err
		}

		if _, err := pool.LookupStorageVolByName(name); err != nil {
			if libvirtErr, ok := err.(libvirt.Error); ok && libvirtErr.Code == libvirt.ERR_NO_STORAGE_VOL {
				break
			}

			return "", "", err
		}

		name = fmt.Sprintf("%s-%d", baseName, i)
		title = fmt.Sprintf("%s-%d", baseTitle, i)
	}

	return name, title, nil
}

func (vmc *VMCreator) createDomainBaseNameAndTitle() (baseName string, baseTitle string) {
	baseTitle = vmc.InstallMedia.Label
	baseName = "unknown"
	return
}

func (vmc *VMCreator) createTargetVolume(name string, storage uint64) (*libvirt.StorageVol, error) {
	pool, err := vmc.EnsureStoragePool()
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

func (vmc *VMCreator) EnsureStoragePool() (*libvirt.StoragePool, error) {
	pool, err := vmc.GetStoragePool()
	if err != nil {
		libvirtErr, ok := err.(libvirt.Error)
		if !ok || libvirtErr.Code != libvirt.ERR_NO_STORAGE_POOL {
			return nil, err
		}

		fmt.Println("Creating storage pool..")
		poolConfig := vmconfigurator.GetPoolConfig()

		xmlConfig, err3 := poolConfig.Marshal()
		if err3 != nil {
			return nil, err3
		}

		pool, err = vmc.connection.StoragePoolDefineXML(xmlConfig, 0)
		if err != nil {
			return nil, err
		}

		if err2 := pool.Build(libvirt.STORAGE_POOL_BUILD_NEW); err2 != nil {
			return nil, err2
		}

		fmt.Println("Created storage pool.")
	}

	// Ensure pool directory exists in case user deleted it after pool creation
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

func (vmc *VMCreator) GetStoragePool() (*libvirt.StoragePool, error) {
	return vmc.connection.LookupStoragePoolByName(config.PackageTarname)
}

func (vmc *VMCreator) RemoveStoragePool() error {
	pool, err := vmc.GetStoragePool()
	if err != nil {
		libvirtErr, ok := err.(libvirt.Error)
		if ok && libvirtErr.Code == libvirt.ERR_NO_STORAGE_POOL {
			return nil
		}

		return err
	}

	if err = pool.Destroy(); err != nil {
		libvirtErr, ok := err.(libvirt.Error)
		if ok && libvirtErr.Code == libvirt.ERR_NO_STORAGE_POOL {
			return nil
		}
	}

	if err = pool.Delete(libvirt.STORAGE_POOL_DELETE_NORMAL); err != nil {
		libvirtErr, ok := err.(libvirt.Error)
		if ok && libvirtErr.Code == libvirt.ERR_NO_STORAGE_POOL {
			return nil
		}
	}

	if err = pool.Undefine(); err != nil {
		libvirtErr, ok := err.(libvirt.Error)
		if ok && libvirtErr.Code == libvirt.ERR_NO_STORAGE_POOL {
			return nil
		}
	}

	return pool.Free()
}
