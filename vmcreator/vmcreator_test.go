package vmcreator_test

import (
	"os"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/require"
	"github.com/tigerinus/libvirt-go-demo/config"
	"github.com/tigerinus/libvirt-go-demo/installermedia"
	"github.com/tigerinus/libvirt-go-demo/vmcreator"
)

func TestCreateVM(t *testing.T) {
	isoPath := "/home/tiger/QEMU/Windows10/Windows10.iso"
	if _, err := os.Stat(isoPath); os.IsNotExist(err) {
		t.Skipf("ISO file not found: %s", isoPath)
	}

	tmpDir, err := os.MkdirTemp("", "libvirt-go-demo-test")
	require.Nil(t, err)

	defer os.RemoveAll(tmpDir)

	xdg.DataHome = tmpDir

	media, err2 := installermedia.ForPath(isoPath)
	require.Nil(t, err2)

	media.Resources.RAM = config.DefaultRAM

	vmc := vmcreator.NewVMCreator(media)
	defer func() {
		err := vmc.RemoveStoragePool()
		require.Nil(t, err)
	}()

	domain, err3 := vmc.CreateVM(false)
	require.Nil(t, err3)

	defer func() {
		err := domain.Destroy()
		require.Nil(t, err)

		err = domain.Undefine()
		require.Nil(t, err)

		err = domain.Free()
		require.Nil(t, err)
	}()

	require.NotNil(t, domain)
}
