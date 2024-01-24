package main

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

func main() {
	fmt.Println("Getting started...")

	conn, err := libvirt.NewConnect("qemu:///session")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	domains, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		panic(err.Error())
	}

	for _, domain := range domains {
		domainName, err := domain.GetName()
		if err != nil {
			_ = fmt.Errorf("error when trying to get the name of a domain: %s", err.Error())
		}

		fmt.Println(domainName)

		if err := domain.Free(); err != nil {
			_ = fmt.Errorf("error when trying to free a domain: %s", err.Error())
		}
	}
}
