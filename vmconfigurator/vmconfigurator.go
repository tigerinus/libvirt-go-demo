package vmconfigurator

import (
	"errors"
	"os"
	"sort"
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
	domain := libvirtxml.Domain{
		Memory: &libvirtxml.DomainMemory{
			Value: media.Resources.RAM,
		},
		Metadata: &libvirtxml.DomainMetadata{
			XML: "<x-casaos>TODO</x-casaos>",
		},
	}

	bestCaps, err := GetBestGuestCaps(caps)
	if err != nil {
		return nil, err
	}

	virtType := "qemu"

	if guestKVMEnabled(*bestCaps) {
		virtType = "kvm"
	}

	setCPUConfig(&domain, caps, &virtType)

	domain.Type = virtType

	setOSConfig(&domain, *bestCaps)

	features := libvirtxml.DomainFeatureList{}

	if guestSupportsFeature(*bestCaps, "acpi") {
		features.ACPI = &libvirtxml.DomainFeature{}
	}

	if guestSupportsFeature(*bestCaps, "apic") {
		features.APIC = &libvirtxml.DomainFeatureAPIC{}
	}

	if guestSupportsFeature(*bestCaps, "pae") {
		features.PAE = &libvirtxml.DomainFeature{}
	}

	domain.Features = &features

	domain.Clock = &libvirtxml.DomainClock{
		Offset: "localtime",
		Timer: []libvirtxml.DomainTimer{
			{
				Name:       "rtc",
				TickPolicy: "catchup",
			},
			{
				Name:       "pit",
				TickPolicy: "delay",
			},
			{
				Name:    "hpet",
				Present: "false",
			},
		},
	}

	return &domain, nil
}

func GetBestGuestCaps(caps libvirtxml.Caps) (*libvirtxml.CapsGuest, error) {
	guestsCaps := caps.Guests

	// Ensure we have the best caps on the top
	sort.Slice(guestsCaps, func(a, b int) bool {
		capsA := guestsCaps[a]
		capsB := guestsCaps[b]

		archA := capsA.Arch.Name
		archB := capsB.Arch.Name

		if archA == "i686" {
			return archB == "x86_64"
		} else if archA == "x86_64" {
			return false
		} else if archB == "x86_64" || archB == "i686" {
			return true
		}

		return false
	})

	for _, guestCaps := range guestsCaps {
		if guestKVMEnabled(guestCaps) {
			return &guestCaps, nil
		}
	}

	for _, guestCaps := range guestsCaps {
		if guestIsQEMU(guestCaps) {
			return &guestCaps, nil
		}
	}

	return nil, errors.New("incapable host system")
}

func setCPUConfig(domain *libvirtxml.Domain, caps libvirtxml.Caps, virtType *string) {
	virtTypeKVM := "kvm"
	if virtType != nil {
		virtType = &virtTypeKVM
	}

	cpuCaps := caps.Host.CPU

	if cpuCaps.Topology == nil {
		return
	}

	topology := &libvirtxml.DomainCPUTopology{
		Sockets: cpuCaps.Topology.Sockets,
		Dies:    cpuCaps.Topology.Dies,
		Cores:   cpuCaps.Topology.Cores,
		Threads: cpuCaps.Topology.Threads,
	}

	domain.VCPU = &libvirtxml.DomainVCPU{
		Value: uint(topology.Sockets * topology.Cores * topology.Threads),
	}

	mode := "host-passthrough"
	if *virtType == "qemu" {
		mode = "host-model"
	}

	domain.CPU = &libvirtxml.DomainCPU{
		Mode:     mode,
		Topology: topology,
	}
}

func setOSConfig(domain *libvirtxml.Domain, guestCaps libvirtxml.CapsGuest) {
	domain.OS = &libvirtxml.DomainOS{
		Type: &libvirtxml.DomainOSType{
			Type: "hvm",
			Arch: guestCaps.Arch.Name,
		},
		// Firmware: "efi",
		BootDevices: []libvirtxml.DomainBootDevice{
			{Dev: "cdrom"},
			{Dev: "hd"},
		},
		BootMenu: &libvirtxml.DomainBootMenu{
			Enable: "true",
		},
	}
}

func guestKVMEnabled(guestCaps libvirtxml.CapsGuest) bool {
	arch := guestCaps.Arch
	for _, domain := range arch.Domains {
		if domain.Type == "kvm" {
			return true
		}
	}

	return false
}

func guestIsQEMU(guestCaps libvirtxml.CapsGuest) bool {
	arch := guestCaps.Arch
	for _, domain := range arch.Domains {
		if domain.Type == "qemu" {
			return true
		}
	}

	return false
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

func guestSupportsFeature(guestCaps libvirtxml.CapsGuest, featureName string) bool {
	if guestCaps.Features == nil {
		return false
	}

	switch featureName {
	case "cpuselection":
		return guestCaps.Features.CPUSelection != nil
	case "deviceboot":
		return guestCaps.Features.DeviceBoot != nil
	case "disksnapshot":
		return guestCaps.Features.DiskSnapshot != nil
	case "pae":
		return guestCaps.Features.PAE != nil
	case "apic":
		return guestCaps.Features.APIC != nil
	case "acpi":
		return guestCaps.Features.ACPI != nil
	case "ia64_be":
		return guestCaps.Features.IA64BE != nil
	}

	return false
}
