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
}
