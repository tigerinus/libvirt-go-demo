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

	defaultPermissions := getDefaultPermissions()

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

func getDefaultPermissions() libvirtxml.StoragePoolTargetPermissions {
	return libvirtxml.StoragePoolTargetPermissions{
		Owner: strconv.Itoa(os.Getuid()),
		Group: strconv.Itoa(os.Getgid()),
		Mode:  "744",
	}
}
