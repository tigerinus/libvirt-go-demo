package vmconfigurator

import (
	"encoding/xml"
	"os"
	"strconv"

	"github.com/tigerinus/libvirt-go-demo/config"
	"github.com/tigerinus/libvirt-go-demo/util"
	"libvirt.org/go/libvirtxml"
)

func GetPoolConfig() (string, error) {
	poolPath := util.GetUserPkgData("images")

	defaultPermissions := getDefaultPoolPermissions()

	config := libvirtxml.StoragePool{
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

	xmlConfig, err := xml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(xmlConfig), nil
}

func CreateVolumeConfig(name string, storage uint64) (string, error) {
	defaultPermissions := getDefaultVolumePermissions()

	config := libvirtxml.StorageVolume{
		Name:     name,
		Capacity: &libvirtxml.StorageVolumeSize{Value: storage},
		Target: &libvirtxml.StorageVolumeTarget{
			Format:      &libvirtxml.StorageVolumeTargetFormat{Type: "qcow2"},
			Compat:      "1.1",
			Permissions: &defaultPermissions,
		},
	}

	xmlConfig, err := xml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(xmlConfig), nil
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
