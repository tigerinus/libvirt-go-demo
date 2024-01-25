package main

import "libvirt.org/go/libvirt"

type VMConfigurator struct{}

func (vmc *VMConfigurator) GetPoolConfig() (*libvirt.StoragePool, error) {
	pool := &libvirt.StoragePool{}
	return pool, nil
}
