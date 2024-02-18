package vmconfigurator

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"

	"github.com/tigerinus/libvirt-go-demo/config"
	"github.com/tigerinus/libvirt-go-demo/installermedia"
	"github.com/tigerinus/libvirt-go-demo/util"
	"libvirt.org/go/libvirtxml"
)

const (
	SpiceAgentChannel = "com.redhat.spice.0"
	WebDAVChannelURI  = "org.spice-space.webdav.0"
)

func CreateDomainConfig(installMedia *installermedia.InstallerMedia, targetPath string, caps libvirtxml.Caps) (*libvirtxml.Domain, error) {
	domain := libvirtxml.Domain{
		Metadata: &libvirtxml.DomainMetadata{
			XML: "<x-casaos>TODO</x-casaos>",
		},
	}

	bestCaps, err := GetBestGuestCaps(caps)
	if err != nil {
		return nil, err
	}

	domain.Memory = &libvirtxml.DomainMemory{
		Unit:  "b",
		Value: installMedia.Resources.RAM,
	}

	virtType := "qemu"

	if guestKVMEnabled(*bestCaps) {
		virtType = "kvm"
	}

	setCPUConfig(&domain, caps, &virtType)

	domain.Type = virtType

	setOSConfig(&domain, installMedia, *bestCaps)

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
				Present: "no",
			},
		},
	}

	SetTargetMediaConfig(&domain, targetPath, installMedia, nil)
	installMedia.SetupDomainConfig(&domain)

	domain.Devices.Graphics = append(domain.Devices.Graphics, libvirtxml.DomainGraphic{
		Spice: CreateGraphicDevice(nil),
	})

	domain.Devices.Channels = append(domain.Devices.Channels, *CreateSpiceAgentChannel())
	domain.Devices.Channels = append(domain.Devices.Channels, *CreateSpiceWebDAVChannel())

	AddUSBSupport(&domain)
	AddSmartcardSupport(&domain)

	setVideoConfig(&domain, installMedia)
	setSoundConfig(&domain, installMedia)
	setTabletConfig(&domain)
	setMouseConfig(&domain)
	setKeyboardConfig(&domain)

	domain.OnPoweroff = "destroy"
	domain.OnReboot = "destroy"
	domain.OnCrash = "destroy"

	// Disable S3 and S4 states for now due to many issues with it currently in qemu/libvirt
	domain.PM = &libvirtxml.DomainPM{
		SuspendToMem: &libvirtxml.DomainPMPolicy{
			Enabled: "no",
		},
		SuspendToDisk: &libvirtxml.DomainPMPolicy{
			Enabled: "no",
		},
	}

	console := libvirtxml.DomainConsole{
		Source: &libvirtxml.DomainChardevSource{
			Pty: &libvirtxml.DomainChardevSourcePty{},
		},
	}

	domain.Devices.Consoles = append(domain.Devices.Consoles, console)

	supportsVirtioNet := installMedia.SupportsVirtIONet() || installMedia.SupportsVirtIO1Net()

	iface := CreateNetworkInterface(&domain, isLibvirtBridgetNetAvailable(), supportsVirtioNet)

	domain.Devices.Interfaces = append(domain.Devices.Interfaces, iface)

	return &domain, nil
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

func GetPoolConfig() *libvirtxml.StoragePool {
	poolPath := util.GetUserPkgData("images")

	defaultPermissions := getDefaultPoolPermissions()

	return &libvirtxml.StoragePool{
		Type:   "dir",
		Name:   config.PackageTarname,
		Source: &libvirtxml.StoragePoolSource{
			// Dir: &libvirtxml.StoragePoolSourceDir{
			// 	Path: poolPath,
			// },
		},
		Target: &libvirtxml.StoragePoolTarget{
			Path:        poolPath,
			Permissions: &defaultPermissions,
		},
	}
}

func SetTargetMediaConfig(domain *libvirtxml.Domain, targetPath string, installMedia *installermedia.InstallerMedia, devIndex *uint8) {
	disk := libvirtxml.DomainDisk{
		// Type: "file",
		Device: "disk",
		Source: &libvirtxml.DomainDiskSource{
			File: &libvirtxml.DomainDiskSourceFile{
				File: targetPath,
			},
		},
	}

	if devIndex == nil {
		devIndex = new(uint8)
		*devIndex = 0
	}

	target := libvirtxml.DomainDiskTarget{}

	driver := libvirtxml.DomainDiskDriver{
		Name:  "qemu",
		Type:  "qcow2",
		Cache: "writeback",
	}

	devLetterStr := string(*devIndex + 97)
	if installMedia.SupportsVirtIODisk() || installMedia.SupportsVirtIO1Disk() {
		fmt.Println("Using virtio controller for the main disk")
		target.Bus = "virtio"
		target.Dev = "vd" + devLetterStr
		driver.Discard = "unmap"
	} else {
		if installMedia.PrefersQ35() {
			fmt.Println("Using SATA controller for the main disk")
			target.Bus = "sata"
			target.Dev = "sd" + devLetterStr
		} else {
			fmt.Println("Using IDE controller for the main disk")
			target.Bus = "ide"
			target.Dev = "hd" + string(97+*devIndex)
		}
	}

	disk.Driver = &driver
	disk.Target = &target

	if domain.Devices == nil {
		domain.Devices = &libvirtxml.DomainDeviceList{}
	}

	domain.Devices.Disks = append(domain.Devices.Disks, disk)
}

func AddSmartcardSupport(domain *libvirtxml.Domain) {
	if domain.Devices == nil {
		domain.Devices = &libvirtxml.DomainDeviceList{}
	}

	domain.Devices.Smartcards = append(domain.Devices.Smartcards, libvirtxml.DomainSmartcard{
		Passthrough: &libvirtxml.DomainChardevSource{
			SpiceVMC: &libvirtxml.DomainChardevSourceSpiceVMC{},
		},
	})
}

func AddUSBSupport(domain *libvirtxml.Domain) {
	if domain.Devices == nil {
		domain.Devices = &libvirtxml.DomainDeviceList{}
	}

	// 4 USB redirection channels
	for i := 0; i < 4; i++ {
		domain.Devices.RedirDevs = append(domain.Devices.RedirDevs, libvirtxml.DomainRedirDev{
			Bus: "usb",
			Source: &libvirtxml.DomainChardevSource{
				SpiceVMC: &libvirtxml.DomainChardevSourceSpiceVMC{},
			},
		})
	}

	controller := createUSBController("qemu-xhci", nil, nil, nil)

	if controller.USB == nil {
		controller.USB = &libvirtxml.DomainControllerUSB{}
	}

	port := uint(15) // 15 is the max amount qemu supports for a single controller
	controller.USB.Port = &port
	domain.Devices.Controllers = append(domain.Devices.Controllers, *controller)
}

func CreateNetworkInterface(domain *libvirtxml.Domain, bridge, virtio bool) libvirtxml.DomainInterface {
	iface := libvirtxml.DomainInterface{
		Source: &libvirtxml.DomainInterfaceSource{},
	}

	if bridge {
		fmt.Printf("Creating bridge network device for %s\n", domain.Name)
		iface.Source.Bridge = &libvirtxml.DomainInterfaceSourceBridge{
			Bridge: "virbr0",
		}
	} else {
		fmt.Printf("Creating user network device for %s\n", domain.Name)
		iface.Source.User = &libvirtxml.DomainInterfaceSourceUser{}
	}

	if virtio {
		iface.Model = &libvirtxml.DomainInterfaceModel{
			Type: "virtio",
		}
	}

	return iface
}

func CreateSpiceAgentChannel() *libvirtxml.DomainChannel {
	return &libvirtxml.DomainChannel{
		Target: &libvirtxml.DomainChannelTarget{
			VirtIO: &libvirtxml.DomainChannelTargetVirtIO{
				Name: SpiceAgentChannel,
			},
		},
		Source: &libvirtxml.DomainChardevSource{
			SpiceVMC: &libvirtxml.DomainChardevSourceSpiceVMC{},
		},
	}
}

func CreateGraphicDevice(accel3d *bool) *libvirtxml.DomainGraphicSpice {
	gl := "no"
	if accel3d != nil && *accel3d {
		gl = "yes"
	}

	return &libvirtxml.DomainGraphicSpice{
		GL: &libvirtxml.DomainGraphicSpiceGL{
			Enable: gl,
		},
		Image: &libvirtxml.DomainGraphicSpiceImage{
			Compression: "off",
		},
	}
}

func CreateSpiceWebDAVChannel() *libvirtxml.DomainChannel {
	return &libvirtxml.DomainChannel{
		Target: &libvirtxml.DomainChannelTarget{
			VirtIO: &libvirtxml.DomainChannelTargetVirtIO{
				Name: WebDAVChannelURI,
			},
		},
		Source: &libvirtxml.DomainChardevSource{
			SpicePort: &libvirtxml.DomainChardevSourceSpicePort{
				Channel: WebDAVChannelURI,
			},
		},
	}
}

func createUSBController(model string, master *libvirtxml.DomainController, index, startPort *uint) *libvirtxml.DomainController {
	if index == nil {
		index = new(uint)
		*index = 0
	}

	if startPort == nil {
		startPort = new(uint)
		*startPort = 0
	}

	controller := libvirtxml.DomainController{
		Model: model,
		Index: index,
		Type:  "usb", // not in GNOME Box
	}

	if master != nil {
		controller.Index = master.Index
		controller.USB = &libvirtxml.DomainControllerUSB{
			Master: &libvirtxml.DomainControllerUSBMaster{
				StartPort: *startPort,
			},
		}
	}

	return &controller
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
			return archB != "x86_64"
		} else if archA == "x86_64" {
			return true
		} else if archB == "x86_64" || archB == "i686" {
			return false
		}

		return true
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

func setOSConfig(domain *libvirtxml.Domain, installMedia *installermedia.InstallerMedia, guestCaps libvirtxml.CapsGuest) {
	osType := libvirtxml.DomainOSType{
		Type: "hvm",
		Arch: guestCaps.Arch.Name,
	}

	if installMedia.PrefersQ35() {
		osType.Machine = "q35" // should it be "pc-q35-8.0"? (see VM under GNOME Box)
	}

	os := libvirtxml.DomainOS{
		Type: &osType,
		// Firmware: "efi",
		BootDevices: []libvirtxml.DomainBootDevice{
			{Dev: "cdrom"},
			{Dev: "hd"},
		},
		BootMenu: &libvirtxml.DomainBootMenu{
			Enable: "yes",
		},
	}

	if installMedia.RequiresEFI() {
		os.Firmware = "efi"
	}

	domain.OS = &os
}

func setVideoConfig(domain *libvirtxml.Domain, installMedia *installermedia.InstallerMedia) {
	video := libvirtxml.DomainVideo{
		Model: libvirtxml.DomainVideoModel{
			Type: "qxl",
		},
	}

	if installMedia.SupportsVirtIOGPU() {
		video.Model.Type = "virtio"
	}

	domain.Devices.Videos = append(domain.Devices.Videos, video)
}

func getSoundModel(installMedia *installermedia.InstallerMedia) string {
	if installMedia.PrefersICH9() {
		return "ich9"
	}

	return "ich6"
}

func setSoundConfig(domain *libvirtxml.Domain, installMedia *installermedia.InstallerMedia) {
	domain.Devices.Sounds = append(domain.Devices.Sounds, libvirtxml.DomainSound{
		Model: getSoundModel(installMedia),
	})
}

func setTabletConfig(domain *libvirtxml.Domain) {
	domain.Devices.Inputs = append(domain.Devices.Inputs, libvirtxml.DomainInput{
		Type: "tablet",
		Bus:  "usb",
	})
}

func setMouseConfig(domain *libvirtxml.Domain) {
	setInputConfig(domain, "mouse")
}

func setKeyboardConfig(domain *libvirtxml.Domain) {
	setInputConfig(domain, "keyboard")
}

func setInputConfig(domain *libvirtxml.Domain, deviceType string) {
	domain.Devices.Inputs = append(domain.Devices.Inputs, libvirtxml.DomainInput{
		Type: deviceType,
		Bus:  "ps2",
	})
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
		Mode:  "0744",
	}
}

func getDefaultVolumePermissions() libvirtxml.StorageVolumeTargetPermissions {
	return libvirtxml.StorageVolumeTargetPermissions{
		Owner: strconv.Itoa(os.Getuid()),
		Group: strconv.Itoa(os.Getgid()),
		Mode:  "0744",
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

func isLibvirtBridgetNetAvailable() bool {
	conn, err := net.Dial("tcp", "localhost:0")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer conn.Close()

	ifName := "virbr0"

	iface, err := net.InterfaceByName(ifName)
	if err != nil || iface.Flags&net.FlagUp == 0 {
		fmt.Printf("Interface '%s' is either not available or not up.\n", ifName)
		return false
	}

	return true
}
