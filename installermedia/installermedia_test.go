package installermedia_test

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tigerinus/libvirt-go-demo/installermedia"
	"libvirt.org/go/libvirtxml"
)

func TestSetupDomainConfig(t *testing.T) {
	domain := libvirtxml.Domain{}

	media, err := installermedia.ForPath("dummy.iso")

	require.Nil(t, err)
	require.NotNil(t, media)

	media.SetupDomainConfig(&domain)

	require.NotNil(t, domain.Devices)
	require.NotNil(t, domain.Devices.Disks)
	require.Len(t, domain.Devices.Disks, 1)

	actual := domain.Devices.Disks[0]

	actualXML, err2 := actual.Marshal()
	require.Nil(t, err2)
	require.NotEmpty(t, actualXML)

	t.Logf("\n%s", actualXML)

	expectedXML := `
	<disk type="file" device="cdrom">
		<driver name="qemu" type="raw"/>
		<source file="dummy.iso" startupPolicy="mandatory"/>
		<target dev="hdc" bus="sata"/>
		<readonly/>
		<address type="drive" controller="0" bus="0" target="0" unit="2"/>
	</disk>  
	`

	var expected libvirtxml.DomainDisk
	err = xml.Unmarshal([]byte(expectedXML), &expected)
	require.Nil(t, err)

	require.Equal(t, expected.Device, actual.Device)
	require.Equal(t, expected.Driver.Name, actual.Driver.Name)
	require.Equal(t, expected.Driver.Type, actual.Driver.Type)
	require.Equal(t, expected.Source.File, actual.Source.File)
	require.Equal(t, expected.Source.StartupPolicy, actual.Source.StartupPolicy)
	require.Equal(t, expected.Target.Dev, actual.Target.Dev)
	require.Equal(t, expected.Target.Bus, actual.Target.Bus)
}
