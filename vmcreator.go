package main

import "libvirt.org/go/libvirt"

type VMCreator struct {
	InstallMedia *InstallerMedia
	connection   *libvirt.Connect
}

func NewVMCreator(installMedia *InstallerMedia) *VMCreator {
	return &VMCreator{
		InstallMedia: installMedia,
	}
}

func (vmc *VMCreator) CreateVM() error {
	name, title, err := vmc.createDomainNameAndTitleFromMedia()
	if err != nil {
		return err
	}

	panic("not implemented" + name + title)
}

func (vmc *VMCreator) createDomainNameAndTitleFromMedia() (name string, title string, err error) {
	baseName, baseTitle := vmc.createDomainBaseNameAndTitle()

	name = baseName
	title = baseTitle

	return
}

func (vmc *VMCreator) createDomainBaseNameAndTitle() (baseName string, baseTitle string) {
	baseTitle = vmc.InstallMedia.Label

	if vmc.InstallMedia.OS != nil {
		baseName = vmc.InstallMedia.OS.ShortID

		if len(baseName) > 12 {
			baseName = baseName[:12]
		}
	} else {
		baseName = "unknown"
	}

	return
}
