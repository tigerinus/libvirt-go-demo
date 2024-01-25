package main

import (
	"fmt"

	"github.com/tigerinus/libvirt-go-demo/config"
	"libvirt.org/go/libvirt"
)

func EnsureStoragePool(connection *libvirt.Connect) error {
	pool, err := GetStoragePool(*connection)
	if err != nil {
		return err
	}

	if pool == nil {
		fmt.Println("Creating storage pool..")

		fmt.Println("Created storage pool.")
	}

	panic("")
}

func GetStoragePool(connection libvirt.Connect) (*libvirt.StoragePool, error) {
	return connection.LookupStoragePoolByName(config.PackageTarname)
}
