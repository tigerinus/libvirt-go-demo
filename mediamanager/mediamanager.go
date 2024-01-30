package mediamanager

import (
	"errors"
	"mime"
	"path/filepath"

	"github.com/tigerinus/libvirt-go-demo/installermedia"
)

type MediaManager struct{}

var mediaManager *MediaManager

var supportedInstallerMediaContentTypes = []string{
	"application/x-cd-image",
}

func pathIsInstallerMedia(path string) bool {
	return MediaMatchesContentType(path, supportedInstallerMediaContentTypes)
}

func MediaMatchesContentType(path string, supportedContentTypes []string) bool {
	contentType := mime.TypeByExtension(filepath.Ext(path))

	for _, supportedContentType := range supportedContentTypes {
		if contentType == supportedContentType {
			return true
		}
	}

	return false
}

func GetDefault() *MediaManager {
	if mediaManager == nil {
		mediaManager = &MediaManager{}
	}

	return mediaManager
}

func (mm *MediaManager) CreateInstallerMediaForPath(path string) (*installermedia.InstallerMedia, error) {
	var media *installermedia.InstallerMedia

	if pathIsInstallerMedia(path) {
		var err error
		media, err = installermedia.ForPath(path)
		if err != nil {
			return nil, err
		}
	}

	if media == nil {
		return nil, errors.New("media is not supported")
	}

	return media, nil
}
