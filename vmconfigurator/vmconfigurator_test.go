package vmconfigurator_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/tigerinus/libvirt-go-demo/installermedia"
	"github.com/tigerinus/libvirt-go-demo/vmconfigurator"
	"libvirt.org/go/libvirtxml"
)

func TestCreateDomainConfig(t *testing.T) {
	t.Fail()
}

func TestCreateVolumeConfig(t *testing.T) {
	actual := vmconfigurator.CreateVolumeConfig("win11", 137438953472)
	require.NotNil(t, actual)

	actualXML, err := actual.Marshal()
	require.Nil(t, err)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<volume type='file'>
			<name>win11</name>
			<capacity unit='bytes'>137438953472</capacity>
			<target>
				<format type='qcow2'/>
				<permissions>
					<mode>0744</mode>
					<owner>1000</owner>
					<group>1000</group>
				</permissions>
				<compat>1.1</compat>
				<path>/home/tiger/.local/share/gnome-boxes/images/win11</path>
				<timestamps>
					<atime>1707285033.566072165</atime>
					<mtime>1707285031.774072176</mtime>
					<ctime>1707285031.774072176</ctime>
					<btime>0</btime>
				</timestamps>
			</target>
			<key>/home/tiger/.local/share/gnome-boxes/images/win11</key>
			<allocation unit='bytes'>200704</allocation>
			<physical unit='bytes'>198656</physical>
		</volume>
	`

	var expected libvirtxml.StorageVolume
	err = expected.Unmarshal(expectedXML)
	require.Nil(t, err)

	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Capacity.Value, actual.Capacity.Value)
	require.Equal(t, expected.Target.Format, actual.Target.Format)
	require.Equal(t, expected.Target.Permissions, actual.Target.Permissions)
	require.Equal(t, expected.Target.Compat, actual.Target.Compat)
}

func TestGetPoolConfig(t *testing.T) {
	actual := vmconfigurator.GetPoolConfig()
	require.NotNil(t, actual)

	actualXML, err := actual.Marshal()
	require.Nil(t, err)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<pool type="dir">
			<name>whatever</name>
			<target>
				<path>whatever/images</path>
				<permissions>
					<owner>1000</owner>
					<group>1000</group>
					<mode>0744</mode>
				</permissions>
			</target>
			<source>
			</source>
			<uuid>886c86cd-7258-4766-98f1-d65e76aa35fd</uuid>
			<capacity unit="bytes">1005867986944</capacity>
			<allocation unit="bytes">112025464832</allocation>
			<available unit="bytes">893842522112</available>
		</pool>
	`

	var expected libvirtxml.StoragePool
	err = expected.Unmarshal(expectedXML)
	require.Nil(t, err)

	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Target.Permissions, actual.Target.Permissions)
}

func TestSetTargetMediaConfig(t *testing.T) {
	domain := libvirtxml.Domain{}

	media, err := installermedia.ForPath("dummy.iso")
	require.Nil(t, err)

	vmconfigurator.SetTargetMediaConfig(&domain, "whatever", media, lo.ToPtr(uint8(0)))

	require.NotNil(t, domain.Devices)
	require.NotNil(t, domain.Devices.Disks)
	require.Len(t, domain.Devices.Disks, 1)

	actual := domain.Devices.Disks[0]
	require.NotNil(t, actual.Target)

	actualXML, err := actual.Marshal()
	require.Nil(t, err)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<disk type="file" device="disk">
			<driver name="qemu" type="qcow2" cache="writeback"/>
			<source file="whatever"/>
			<target dev="sda" bus="sata"/>
			<address type="drive" controller="0" bus="0" target="0" unit="0"/>
		</disk>
	`

	var expected libvirtxml.DomainDisk
	err = expected.Unmarshal(expectedXML)
	require.Nil(t, err)

	require.Equal(t, expected.Device, actual.Device)
	require.Equal(t, expected.Driver, actual.Driver)
	require.Equal(t, expected.Source, actual.Source)
	require.Equal(t, expected.Target, actual.Target)
}

func TestAddSmartcardSupport(t *testing.T) {
	domain := libvirtxml.Domain{}

	vmconfigurator.AddSmartcardSupport(&domain)

	require.NotNil(t, domain.Devices)
	require.NotNil(t, domain.Devices.Smartcards)
	require.Len(t, domain.Devices.Smartcards, 1)

	actual := domain.Devices.Smartcards[0]
	require.NotNil(t, actual)

	actualXML, err := actual.Marshal()
	require.Nil(t, err)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<smartcard mode="passthrough" type="spicevmc">
			<address type="ccid" controller="0" slot="0"/>
		</smartcard>  
	`

	var expected libvirtxml.DomainSmartcard
	err = expected.Unmarshal(expectedXML)
	require.Nil(t, err)

	require.Equal(t, expected.Passthrough, actual.Passthrough)
}
