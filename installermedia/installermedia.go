package installermedia

import (
	"fmt"
	"path/filepath"
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

func (im *InstallerMedia) labelSetup(label *string) {
	if label != nil {
		im.Label = *label
		return
	}

	im.Label = filepath.Base(im.DeviceFile)
}
