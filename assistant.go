package main

import (
	"fmt"

	"github.com/tigerinus/libvirt-go-demo/installermedia"
	"github.com/tigerinus/libvirt-go-demo/mediamanager"
)

type Assistant struct {
	installerMedia *installermedia.InstallerMedia
}

func (a *Assistant) PrepareForPath(path string) {
	a.installerMedia = createInstallerMedia(path)
}

func createInstallerMedia(path string) *installermedia.InstallerMedia {
	mediaManager := mediamanager.GetDefault()

	var media *installermedia.InstallerMedia
	var err error

	media, err = mediaManager.CreateInstallerMediaForPath(path)

	if err != nil {
		fmt.Println(err.Error())
	}

	return media
}
