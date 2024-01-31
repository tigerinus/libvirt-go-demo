package vmconfigurator

import (
	"os"
	"strconv"

	"github.com/tigerinus/libvirt-go-demo/config"
	"github.com/tigerinus/libvirt-go-demo/installermedia"
	"github.com/tigerinus/libvirt-go-demo/util"
	"libvirt.org/go/libvirtxml"
)

func GetPoolConfig() *libvirtxml.StoragePool {
	poolPath := util.GetUserPkgData("images")

	defaultPermissions := getDefaultPoolPermissions()

	return &libvirtxml.StoragePool{
		Type: "dir",
		Name: config.PackageTarname,
		Source: &libvirtxml.StoragePoolSource{
			Dir: &libvirtxml.StoragePoolSourceDir{
				Path: poolPath,
			},
		},
		Target: &libvirtxml.StoragePoolTarget{
			Path:        poolPath,
			Permissions: &defaultPermissions,
		},
	}
}

func CreateVolumeConfig(name string, storage uint64) *libvirtxml.StorageVolume {
	defaultPermissions := getDefaultVolumePermissions()

	return &libvirtxml.StorageVolume{
		Name:     name,
		Capacity: &libvirtxml.StorageVolumeSize{Value: storage},
		Target: &libvirtxml.StorageVolumeTarget{
			Format:      &libvirtxml.StorageVolumeTargetFormat{Type: "qcow2"},
			Compat:      "1.1",
			Permissions: &defaultPermissions,
		},
	}
}

func CreateDomainConfig(media *installermedia.InstallerMedia, targetPath string, caps libvirtxml.Caps, domainCaps libvirtxml.DomainCaps) (*libvirtxml.Domain, error) {
	config := libvirtxml.Domain{}

	return &config, nil
}

func getDefaultPoolPermissions() libvirtxml.StoragePoolTargetPermissions {
	return libvirtxml.StoragePoolTargetPermissions{
		Owner: strconv.Itoa(os.Getuid()),
		Group: strconv.Itoa(os.Getgid()),
		Mode:  "744",
	}
}

func getDefaultVolumePermissions() libvirtxml.StorageVolumeTargetPermissions {
	return libvirtxml.StorageVolumeTargetPermissions{
		Owner: strconv.Itoa(os.Getuid()),
		Group: strconv.Itoa(os.Getgid()),
		Mode:  "744",
	}
}
