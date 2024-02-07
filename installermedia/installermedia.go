package installermedia

import (
	"fmt"
	"path/filepath"

	"libvirt.org/go/libvirtxml"
)

type Resources struct {
	RAM uint
}

type InstallerMedia struct {
	Label      string
	DeviceFile string
	FromImage  bool
	SkipImport bool
	Resources  Resources
}

func ForPath(path string) (*InstallerMedia, error) {
	media := InstallerMedia{
		DeviceFile: path,
		FromImage:  true,
	}

	media.labelSetup(nil)

	return &media, nil
}

// dummy func for now
func (im *InstallerMedia) PrepareForInstallation(vmName string) {
	fmt.Printf("PrepareForInstallation(\"%s\") is not implemented for now\n", vmName)
}

func (im *InstallerMedia) SetupDomainConfig(domain *libvirtxml.Domain) {
	diskType := "block"
	if im.FromImage {
		diskType = "file"
	}

	mandatory := true

	im.addCDConfig(domain, diskType, &im.DeviceFile, "hdc", &mandatory)
}

func (im *InstallerMedia) addCDConfig(
	domain *libvirtxml.Domain,
	diskType string,
	isoPath *string,
	deviceName string,
	mandatory *bool,
) {
	disk := libvirtxml.DomainDisk{
		// Type: diskType,
		Device: "cdrom",
		Driver: &libvirtxml.DomainDiskDriver{
			Name: "qemu",
			Type: "raw",
		},
		Target: &libvirtxml.DomainDiskTarget{
			Dev: deviceName,
		},
	}

	if isoPath != nil {
		disk.Source = &libvirtxml.DomainDiskSource{
			File: &libvirtxml.DomainDiskSourceFile{
				File: *isoPath,
			},
		}
	}

	if im.PrefersQ35() {
		disk.Target.Bus = "sata"
	} else {
		disk.Target.Bus = "ide"
	}

	if mandatory == nil {
		mandatory = new(bool)
		*mandatory = false
	}

	if diskType == "file" && *mandatory {
		disk.Source.StartupPolicy = "mandatory"
	}

	if domain.Devices == nil {
		domain.Devices = &libvirtxml.DomainDeviceList{}
	}

	domain.Devices.Disks = append(domain.Devices.Disks, disk)
}

func (im *InstallerMedia) labelSetup(label *string) {
	if label != nil {
		im.Label = *label
		return
	}

	im.Label = filepath.Base(im.DeviceFile)
}

// getters

func (im *InstallerMedia) SupportsVirtIO1Disk() bool {
	// TEMP: for now, we'll just return false
	return false
}

func (im *InstallerMedia) SupportsVirtIODisk() bool {
	// TEMP: for now, we'll just return false
	return false
}

func (im *InstallerMedia) SupportsVirtIO1Net() bool {
	// TEMP: for now, we'll just return false
	return false
}

func (im *InstallerMedia) SupportsVirtIONet() bool {
	// TEMP: for now, we'll just return false
	return false
}

func (im *InstallerMedia) SupportsVirtIOGPU() bool {
	// TEMP: for now, we'll just return false
	return false
}

func (im *InstallerMedia) PrefersQ35() bool {
	if im.SupportsVirtIONet() && !im.SupportsVirtIO1Net() {
		return false
	}

	// TEMP: for now, we'll just return true
	return true
}

func (im *InstallerMedia) PrefersICH9() bool {
	// TEMP: for now, we'll just return the same value as PrefersQ35()
	return im.PrefersQ35()
}

func (im *InstallerMedia) RequiresEFI() bool {
	// TEMP: for now, we'll just return false
	return false
}
