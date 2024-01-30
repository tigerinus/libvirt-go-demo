package installermedia

import "path/filepath"

type InstallerMedia struct {
	Label      string
	DeviceFile string
	FromImage  bool
}

func ForPath(path string) (*InstallerMedia, error) {
	media := InstallerMedia{
		DeviceFile: path,
		FromImage:  true,
	}

	media.labelSetup(nil)

	return &media, nil
}

func (im *InstallerMedia) labelSetup(label *string) {
	if label != nil {
		im.Label = *label
		return
	}

	im.Label = filepath.Base(im.DeviceFile)
}
