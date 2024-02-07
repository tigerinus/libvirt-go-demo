package vmconfigurator_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tigerinus/libvirt-go-demo/vmconfigurator"
	"libvirt.org/go/libvirtxml"
)

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
