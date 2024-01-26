package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/tigerinus/libvirt-go-demo/config"
	"libvirt.org/go/libvirt"
)

func EnsureStoragePool(connection *libvirt.Connect) (*libvirt.StoragePool, error) {
	pool, err := GetStoragePool(*connection)
	if err != nil {
		return nil, err
	}

	if pool == nil {
		fmt.Println("Creating storage pool..")
		poolConfig, err2 := GetPoolConfig()
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

	poolPath := GetUserPkgData("images")
	EnsureDirectory(poolPath)

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

func GetUserPkgData(fileName string) string {
	dir := filepath.Join(xdg.DataHome, config.PackageTarname)

	EnsureDirectory(dir)

	return filepath.Join(dir, fileName)
}

func EnsureDirectory(dir string) {
	if err := os.MkdirAll(dir, 0x744); err != nil {
		fmt.Printf("warning: %s", err.Error())
	}
}
