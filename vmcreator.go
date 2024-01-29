package main

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

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

	pool, err := EnsureStoragePool(vmc.connection)

	for i := 2; ; i++ {
		domain, err := vmc.connection.LookupDomainByName(name)
		if err != nil {
			return "", "", err
		}

		if domain == nil {
			break
		}

		domain, err = vmc.connection.LookupDomainByName(title)
		if err != nil {
			return "", "", err
		}

		if domain == nil {
			break
		}

		volume, err := pool.LookupStorageVolByName(name)
		if err != nil {
			return "", "", err
		}

		if volume == nil {
			break
		}

		name = fmt.Sprintf("%s-%d", baseName, i)
		title = fmt.Sprintf("%s-%d", baseTitle, i)
	}

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
