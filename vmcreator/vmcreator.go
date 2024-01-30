package vmcreator

import (
	"fmt"

	"github.com/tigerinus/libvirt-go-demo/config"
	"github.com/tigerinus/libvirt-go-demo/installermedia"
	"github.com/tigerinus/libvirt-go-demo/util"
	"github.com/tigerinus/libvirt-go-demo/vmconfigurator"
	"libvirt.org/go/libvirt"
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

func (vmc *VMCreator) CreateVM() error {
	name, title, err := vmc.createDomainNameAndTitleFromMedia()
	if err != nil {
		return err
	}

	panic("not implemented" + name + title)
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

func EnsureStoragePool(connection *libvirt.Connect) (*libvirt.StoragePool, error) {
	pool, err := GetStoragePool(*connection)
	if err != nil {
		return nil, err
	}

	if pool == nil {
		fmt.Println("Creating storage pool..")
		poolConfig, err2 := vmconfigurator.GetPoolConfig()
		if err2 != nil {
			return nil, err2
		}

		pool, err2 = connection.StoragePoolCreateXML(poolConfig, libvirt.STORAGE_POOL_CREATE_NORMAL)
		if err2 != nil {
			return nil, err2
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