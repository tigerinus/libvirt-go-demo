package main

type InstallerMedia struct {
	OS       *OSInfo
	Label    string
	DiskFile string

	hostname string
}

func (im *InstallerMedia) PrepareForInstallation(vmName string) {
	im.PrepareToContinueInstallation(vmName)

	// TODO
}

func (im *InstallerMedia) PrepareToContinueInstallation(vmName string) {
	im.hostname = ReplaceRegex(vmName, "[{|}~[\\]^':; <=>?@!\"#$%`()+/.,*&]", "")

	unattended := "unattended.iso"

	path := im.GetUserUnattended(&unattended)

	im.DiskFile = path
}

func (im *InstallerMedia) GetUserUnattended(suffix *string) string {
	filename := im.hostname
	if suffix != nil {
		filename += "-" + *suffix
	}

	return GetUserPkgCache(filename)
}
