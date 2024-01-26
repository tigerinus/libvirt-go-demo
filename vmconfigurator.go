package main

import (
	"encoding/xml"
	"os"
	"strconv"

	"github.com/tigerinus/libvirt-go-demo/config"
	"github.com/tigerinus/libvirt-go-demo/schema"
)

func GetPoolConfig() (string, error) {
	poolPath := GetUserPkgData("images")

	config := schema.StoragePool{
		Type: schema.StoragePoolTypeDir,
		Name: config.PackageTarname,
		Source: schema.StoragePoolSource{
			Directory: schema.StoragePoolSourceDirectory{
				Path: poolPath,
			},
		},
		Target: schema.StoragePoolTarget{
			Path:        poolPath,
			Permissions: getDefaultPermissions(),
		},
	}

	xmlConfig, err := xml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(xmlConfig), nil
}

func getDefaultPermissions() schema.StoragePoolTargetPermissions {
	return schema.StoragePoolTargetPermissions{
		Owner: strconv.Itoa(os.Getuid()),
		Group: strconv.Itoa(os.Getgid()),
		Mode:  "744",
	}
}
