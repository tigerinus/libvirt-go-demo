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
}

func (im *InstallerMedia) addCDConfig(
	domain *libvirtxml.Domain,
	diskType libvirtxml.DomainDisk,
) {
	disk := libvirtxml.DomainDisk{}
}

func (im *InstallerMedia) labelSetup(label *string) {
	if label != nil {
		im.Label = *label
		return
	}

	im.Label = filepath.Base(im.DeviceFile)
}

// getters
func (im *InstallerMedia) SupportsVirtIODisk() bool {
	// TEMP: for now, we'll just return false
	return false
}

func (im *InstallerMedia) SupportsVirtIO1Disk() bool {
	// TEMP: for now, we'll just return false
	return false
}

func (im *InstallerMedia) SupportsVirtIONet() bool {
	// TEMP: for now, we'll just return false
	return false
}

func (im *InstallerMedia) SupportsVirtIO1Net() bool {
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

func (im *InstallerMedia) RequiresEFI() bool {
	// TEMP: for now, we'll just return false
	return false
}
