package vmconfigurator_test

import (
	"encoding/xml"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/tigerinus/libvirt-go-demo/installermedia"
	"github.com/tigerinus/libvirt-go-demo/vmconfigurator"
	"libvirt.org/go/libvirtxml"
)

func TestCreateDomainConfig(t *testing.T) {
	capsXML := `
	<capabilities>
		<host>
		<uuid>a85770b3-3824-4357-8058-4fea9d35e862</uuid>
		<cpu>
			<arch>x86_64</arch>
			<model>Broadwell-noTSX-IBRS</model>
			<vendor>Intel</vendor>
			<microcode version='1074'/>
			<signature family='6' model='154' stepping='3'/>
			<counter name='tsc' frequency='3110400000' scaling='yes'/>
			<topology sockets='1' dies='1' cores='16' threads='1'/>
			<maxphysaddr mode='emulate' bits='39'/>
			<feature name='vme'/>
			<feature name='ds'/>
			<feature name='acpi'/>
			<feature name='ss'/>
			<feature name='ht'/>
			<feature name='tm'/>
			<feature name='pbe'/>
			<feature name='dtes64'/>
			<feature name='monitor'/>
			<feature name='ds_cpl'/>
			<feature name='vmx'/>
			<feature name='smx'/>
			<feature name='est'/>
			<feature name='tm2'/>
			<feature name='xtpr'/>
			<feature name='pdcm'/>
			<feature name='osxsave'/>
			<feature name='f16c'/>
			<feature name='rdrand'/>
			<feature name='arat'/>
			<feature name='tsc_adjust'/>
			<feature name='clflushopt'/>
			<feature name='clwb'/>
			<feature name='intel-pt'/>
			<feature name='sha-ni'/>
			<feature name='umip'/>
			<feature name='pku'/>
			<feature name='ospke'/>
			<feature name='waitpkg'/>
			<feature name='gfni'/>
			<feature name='vaes'/>
			<feature name='vpclmulqdq'/>
			<feature name='rdpid'/>
			<feature name='movdiri'/>
			<feature name='movdir64b'/>
			<feature name='pks'/>
			<feature name='fsrm'/>
			<feature name='md-clear'/>
			<feature name='serialize'/>
			<feature name='arch-lbr'/>
			<feature name='stibp'/>
			<feature name='flush-l1d'/>
			<feature name='arch-capabilities'/>
			<feature name='core-capability'/>
			<feature name='ssbd'/>
			<feature name='avx-vnni'/>
			<feature name='fsrs'/>
			<feature name='xsaveopt'/>
			<feature name='xsavec'/>
			<feature name='xgetbv1'/>
			<feature name='xsaves'/>
			<feature name='pdpe1gb'/>
			<feature name='abm'/>
			<feature name='invtsc'/>
			<feature name='rdctl-no'/>
			<feature name='ibrs-all'/>
			<feature name='skip-l1dfl-vmentry'/>
			<feature name='mds-no'/>
			<feature name='pschange-mc-no'/>
			<feature name='sbdr-ssdp-no'/>
			<feature name='fbsdp-no'/>
			<feature name='psdp-no'/>
			<pages unit='KiB' size='4'/>
			<pages unit='KiB' size='2048'/>
			<pages unit='KiB' size='1048576'/>
		</cpu>
		<power_management>
			<suspend_mem/>
		</power_management>
		<iommu support='yes'/>
		<migration_features>
			<live/>
			<uri_transports>
			<uri_transport>tcp</uri_transport>
			<uri_transport>rdma</uri_transport>
			</uri_transports>
		</migration_features>
		<topology>
			<cells num='1'>
			<cell id='0'>
				<memory unit='KiB'>16084312</memory>
				<pages unit='KiB' size='4'>4021078</pages>
				<pages unit='KiB' size='2048'>0</pages>
				<pages unit='KiB' size='1048576'>0</pages>
				<distances>
				<sibling id='0' value='10'/>
				</distances>
				<cpus num='16'>
				<cpu id='0' socket_id='0' die_id='0' core_id='0' siblings='0-1'/>
				<cpu id='1' socket_id='0' die_id='0' core_id='0' siblings='0-1'/>
				<cpu id='2' socket_id='0' die_id='0' core_id='4' siblings='2-3'/>
				<cpu id='3' socket_id='0' die_id='0' core_id='4' siblings='2-3'/>
				<cpu id='4' socket_id='0' die_id='0' core_id='8' siblings='4-5'/>
				<cpu id='5' socket_id='0' die_id='0' core_id='8' siblings='4-5'/>
				<cpu id='6' socket_id='0' die_id='0' core_id='12' siblings='6-7'/>
				<cpu id='7' socket_id='0' die_id='0' core_id='12' siblings='6-7'/>
				<cpu id='8' socket_id='0' die_id='0' core_id='16' siblings='8'/>
				<cpu id='9' socket_id='0' die_id='0' core_id='17' siblings='9'/>
				<cpu id='10' socket_id='0' die_id='0' core_id='18' siblings='10'/>
				<cpu id='11' socket_id='0' die_id='0' core_id='19' siblings='11'/>
				<cpu id='12' socket_id='0' die_id='0' core_id='20' siblings='12'/>
				<cpu id='13' socket_id='0' die_id='0' core_id='21' siblings='13'/>
				<cpu id='14' socket_id='0' die_id='0' core_id='22' siblings='14'/>
				<cpu id='15' socket_id='0' die_id='0' core_id='23' siblings='15'/>
				</cpus>
			</cell>
			</cells>
		</topology>
		<cache>
			<bank id='0' level='3' type='both' size='18' unit='MiB' cpus='0-15'/>
		</cache>
		<secmodel>
			<model>none</model>
			<doi>0</doi>
		</secmodel>
		</host>
	
		<guest>
		<os_type>hvm</os_type>
		<arch name='i686'>
			<wordsize>32</wordsize>
			<emulator>/usr/bin/qemu-system-i386</emulator>
			<machine maxCpus='255'>pc-i440fx-mantic</machine>
			<machine canonical='pc-i440fx-mantic' maxCpus='255'>ubuntu</machine>
			<machine maxCpus='288'>pc-q35-mantic</machine>
			<machine canonical='pc-q35-mantic' maxCpus='288'>ubuntu-q35</machine>
			<machine maxCpus='255'>pc-i440fx-impish-hpb</machine>
			<machine maxCpus='288'>pc-q35-5.2</machine>
			<machine maxCpus='288'>pc-q35-lunar-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-2.12</machine>
			<machine maxCpus='255'>pc-i440fx-2.0</machine>
			<machine maxCpus='255'>pc-i440fx-xenial</machine>
			<machine maxCpus='288'>pc-q35-kinetic</machine>
			<machine maxCpus='255'>pc-i440fx-6.2</machine>
			<machine maxCpus='288'>pc-q35-4.2</machine>
			<machine maxCpus='255'>pc-i440fx-2.5</machine>
			<machine maxCpus='255'>pc-i440fx-4.2</machine>
			<machine maxCpus='255'>pc-i440fx-focal</machine>
			<machine maxCpus='255'>pc-i440fx-hirsute</machine>
			<machine maxCpus='255'>pc-q35-xenial</machine>
			<machine maxCpus='255'>pc-i440fx-jammy-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-5.2</machine>
			<machine maxCpus='255' deprecated='yes'>pc-i440fx-1.5</machine>
			<machine maxCpus='255'>pc-q35-2.7</machine>
			<machine maxCpus='288'>pc-q35-eoan-hpb</machine>
			<machine maxCpus='288'>pc-q35-groovy</machine>
			<machine maxCpus='255'>pc-i440fx-disco-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-zesty</machine>
			<machine maxCpus='288'>pc-q35-lunar</machine>
			<machine maxCpus='255'>pc-i440fx-groovy</machine>
			<machine maxCpus='288'>pc-q35-7.1</machine>
			<machine maxCpus='288'>pc-q35-artful</machine>
			<machine maxCpus='255'>pc-i440fx-2.2</machine>
			<machine maxCpus='255'>pc-i440fx-trusty</machine>
			<machine maxCpus='255'>pc-i440fx-eoan-hpb</machine>
			<machine maxCpus='288'>pc-q35-focal-hpb</machine>
			<machine maxCpus='288'>pc-q35-bionic-hpb</machine>
			<machine maxCpus='288'>pc-q35-mantic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-artful</machine>
			<machine maxCpus='255'>pc-i440fx-2.7</machine>
			<machine maxCpus='288'>pc-q35-6.1</machine>
			<machine maxCpus='255'>pc-i440fx-kinetic</machine>
			<machine maxCpus='255'>pc-q35-2.4</machine>
			<machine maxCpus='255'>pc-i440fx-yakkety</machine>
			<machine maxCpus='288'>pc-q35-cosmic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-7.1</machine>
			<machine maxCpus='288'>pc-q35-2.10</machine>
			<machine maxCpus='1'>x-remote</machine>
			<machine maxCpus='255' deprecated='yes'>pc-i440fx-1.7</machine>
			<machine maxCpus='288'>pc-q35-5.1</machine>
			<machine maxCpus='288'>pc-q35-2.9</machine>
			<machine maxCpus='255'>pc-i440fx-2.11</machine>
			<machine maxCpus='288'>pc-q35-3.1</machine>
			<machine maxCpus='255'>pc-i440fx-6.1</machine>
			<machine maxCpus='288'>pc-q35-4.1</machine>
			<machine maxCpus='288'>pc-q35-jammy</machine>
			<machine maxCpus='255'>pc-i440fx-2.4</machine>
			<machine maxCpus='255'>pc-i440fx-4.1</machine>
			<machine maxCpus='288'>pc-q35-eoan</machine>
			<machine maxCpus='288'>pc-q35-jammy-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-5.1</machine>
			<machine maxCpus='255'>pc-i440fx-2.9</machine>
			<machine maxCpus='255'>pc-i440fx-bionic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-lunar</machine>
			<machine maxCpus='1'>isapc</machine>
			<machine maxCpus='255'>pc-i440fx-mantic-hpb</machine>
			<machine maxCpus='255' deprecated='yes'>pc-i440fx-1.4</machine>
			<machine maxCpus='288'>pc-q35-cosmic</machine>
			<machine maxCpus='255'>pc-q35-2.6</machine>
			<machine maxCpus='255'>pc-i440fx-3.1</machine>
			<machine maxCpus='288'>pc-q35-bionic</machine>
			<machine maxCpus='288'>pc-q35-disco-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-cosmic</machine>
			<machine maxCpus='288'>pc-q35-2.12</machine>
			<machine maxCpus='255'>pc-i440fx-bionic</machine>
			<machine maxCpus='288'>pc-q35-kinetic-hpb</machine>
			<machine maxCpus='288'>pc-q35-7.0</machine>
			<machine maxCpus='255'>pc-i440fx-lunar-hpb</machine>
			<machine maxCpus='288'>pc-q35-disco</machine>
			<machine maxCpus='255'>pc-i440fx-cosmic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-2.1</machine>
			<machine maxCpus='288'>pc-q35-groovy-hpb</machine>
			<machine maxCpus='288'>pc-q35-8.0</machine>
			<machine canonical='pc-q35-8.0' maxCpus='288'>q35</machine>
			<machine maxCpus='288'>pc-q35-impish</machine>
			<machine maxCpus='255'>pc-i440fx-wily</machine>
			<machine maxCpus='255'>pc-i440fx-8.0</machine>
			<machine canonical='pc-i440fx-8.0' maxCpus='255'>pc</machine>
			<machine maxCpus='288'>pc-q35-6.0</machine>
			<machine maxCpus='255'>pc-i440fx-impish</machine>
			<machine maxCpus='255'>pc-i440fx-2.6</machine>
			<machine maxCpus='255'>pc-i440fx-jammy</machine>
			<machine maxCpus='288'>pc-q35-impish-hpb</machine>
			<machine maxCpus='288'>pc-q35-hirsute</machine>
			<machine maxCpus='288'>pc-q35-4.0.1</machine>
			<machine maxCpus='288'>pc-q35-hirsute-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-7.0</machine>
			<machine maxCpus='255' deprecated='yes'>pc-i440fx-1.6</machine>
			<machine maxCpus='288'>pc-q35-5.0</machine>
			<machine maxCpus='288'>pc-q35-2.8</machine>
			<machine maxCpus='255'>pc-i440fx-2.10</machine>
			<machine maxCpus='288'>pc-q35-3.0</machine>
			<machine maxCpus='255'>pc-i440fx-6.0</machine>
			<machine maxCpus='288'>pc-q35-zesty</machine>
			<machine maxCpus='288'>pc-q35-7.2</machine>
			<machine maxCpus='288'>pc-q35-4.0</machine>
			<machine maxCpus='288'>pc-q35-focal</machine>
			<machine maxCpus='288'>microvm</machine>
			<machine maxCpus='255'>pc-i440fx-2.3</machine>
			<machine maxCpus='255'>pc-i440fx-disco</machine>
			<machine maxCpus='255'>pc-i440fx-kinetic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-4.0</machine>
			<machine maxCpus='255'>pc-i440fx-focal-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-groovy-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-hirsute-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-5.0</machine>
			<machine maxCpus='255'>pc-i440fx-2.8</machine>
			<machine maxCpus='288'>pc-q35-6.2</machine>
			<machine maxCpus='255'>pc-i440fx-eoan</machine>
			<machine maxCpus='255'>pc-q35-2.5</machine>
			<machine maxCpus='255'>pc-i440fx-3.0</machine>
			<machine maxCpus='255'>pc-q35-yakkety</machine>
			<machine maxCpus='255'>pc-i440fx-7.2</machine>
			<machine maxCpus='288'>pc-q35-2.11</machine>
			<domain type='qemu'/>
			<domain type='kvm'/>
		</arch>
		<features>
			<pae/>
			<nonpae/>
			<acpi default='on' toggle='yes'/>
			<apic default='on' toggle='no'/>
			<cpuselection/>
			<deviceboot/>
			<disksnapshot default='on' toggle='no'/>
		</features>
		</guest>
	
		<guest>
		<os_type>hvm</os_type>
		<arch name='x86_64'>
			<wordsize>64</wordsize>
			<emulator>/usr/bin/qemu-system-x86_64</emulator>
			<machine maxCpus='255'>pc-i440fx-mantic</machine>
			<machine canonical='pc-i440fx-mantic' maxCpus='255'>ubuntu</machine>
			<machine maxCpus='288'>pc-q35-mantic</machine>
			<machine canonical='pc-q35-mantic' maxCpus='288'>ubuntu-q35</machine>
			<machine maxCpus='255'>pc-i440fx-impish-hpb</machine>
			<machine maxCpus='288'>pc-q35-5.2</machine>
			<machine maxCpus='288'>pc-q35-lunar-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-2.12</machine>
			<machine maxCpus='255'>pc-i440fx-2.0</machine>
			<machine maxCpus='255'>pc-i440fx-xenial</machine>
			<machine maxCpus='288'>pc-q35-kinetic</machine>
			<machine maxCpus='255'>pc-i440fx-6.2</machine>
			<machine maxCpus='288'>pc-q35-4.2</machine>
			<machine maxCpus='255'>pc-i440fx-2.5</machine>
			<machine maxCpus='255'>pc-i440fx-4.2</machine>
			<machine maxCpus='255'>pc-i440fx-hirsute</machine>
			<machine maxCpus='255'>pc-i440fx-focal</machine>
			<machine maxCpus='255'>pc-q35-xenial</machine>
			<machine maxCpus='255'>pc-i440fx-jammy-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-5.2</machine>
			<machine maxCpus='255' deprecated='yes'>pc-i440fx-1.5</machine>
			<machine maxCpus='255'>pc-q35-2.7</machine>
			<machine maxCpus='288'>pc-q35-eoan-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-zesty</machine>
			<machine maxCpus='255'>pc-i440fx-disco-hpb</machine>
			<machine maxCpus='288'>pc-q35-groovy</machine>
			<machine maxCpus='288'>pc-q35-lunar</machine>
			<machine maxCpus='255'>pc-i440fx-groovy</machine>
			<machine maxCpus='288'>pc-q35-7.1</machine>
			<machine maxCpus='288'>pc-q35-artful</machine>
			<machine maxCpus='255'>pc-i440fx-trusty</machine>
			<machine maxCpus='255'>pc-i440fx-2.2</machine>
			<machine maxCpus='288'>pc-q35-focal-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-eoan-hpb</machine>
			<machine maxCpus='288'>pc-q35-bionic-hpb</machine>
			<machine maxCpus='288'>pc-q35-mantic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-artful</machine>
			<machine maxCpus='255'>pc-i440fx-2.7</machine>
			<machine maxCpus='288'>pc-q35-6.1</machine>
			<machine maxCpus='255'>pc-i440fx-kinetic</machine>
			<machine maxCpus='255'>pc-i440fx-yakkety</machine>
			<machine maxCpus='255'>pc-q35-2.4</machine>
			<machine maxCpus='288'>pc-q35-cosmic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-7.1</machine>
			<machine maxCpus='288'>pc-q35-2.10</machine>
			<machine maxCpus='1'>x-remote</machine>
			<machine maxCpus='288'>pc-q35-5.1</machine>
			<machine maxCpus='255' deprecated='yes'>pc-i440fx-1.7</machine>
			<machine maxCpus='288'>pc-q35-2.9</machine>
			<machine maxCpus='255'>pc-i440fx-2.11</machine>
			<machine maxCpus='288'>pc-q35-3.1</machine>
			<machine maxCpus='255'>pc-i440fx-6.1</machine>
			<machine maxCpus='288'>pc-q35-4.1</machine>
			<machine maxCpus='288'>pc-q35-jammy</machine>
			<machine maxCpus='255'>pc-i440fx-2.4</machine>
			<machine maxCpus='255'>pc-i440fx-4.1</machine>
			<machine maxCpus='288'>pc-q35-eoan</machine>
			<machine maxCpus='288'>pc-q35-jammy-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-5.1</machine>
			<machine maxCpus='255'>pc-i440fx-2.9</machine>
			<machine maxCpus='255'>pc-i440fx-bionic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-lunar</machine>
			<machine maxCpus='1'>isapc</machine>
			<machine maxCpus='255'>pc-i440fx-mantic-hpb</machine>
			<machine maxCpus='255' deprecated='yes'>pc-i440fx-1.4</machine>
			<machine maxCpus='288'>pc-q35-cosmic</machine>
			<machine maxCpus='255'>pc-q35-2.6</machine>
			<machine maxCpus='255'>pc-i440fx-3.1</machine>
			<machine maxCpus='288'>pc-q35-bionic</machine>
			<machine maxCpus='288'>pc-q35-disco-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-cosmic</machine>
			<machine maxCpus='288'>pc-q35-2.12</machine>
			<machine maxCpus='255'>pc-i440fx-bionic</machine>
			<machine maxCpus='288'>pc-q35-kinetic-hpb</machine>
			<machine maxCpus='288'>pc-q35-groovy-hpb</machine>
			<machine maxCpus='288'>pc-q35-7.0</machine>
			<machine maxCpus='255'>pc-i440fx-lunar-hpb</machine>
			<machine maxCpus='288'>pc-q35-disco</machine>
			<machine maxCpus='255'>pc-i440fx-cosmic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-2.1</machine>
			<machine maxCpus='288'>pc-q35-8.0</machine>
			<machine canonical='pc-q35-8.0' maxCpus='288'>q35</machine>
			<machine maxCpus='255'>pc-i440fx-wily</machine>
			<machine maxCpus='288'>pc-q35-impish</machine>
			<machine maxCpus='255'>pc-i440fx-8.0</machine>
			<machine canonical='pc-i440fx-8.0' maxCpus='255'>pc</machine>
			<machine maxCpus='255'>pc-i440fx-2.6</machine>
			<machine maxCpus='288'>pc-q35-6.0</machine>
			<machine maxCpus='255'>pc-i440fx-impish</machine>
			<machine maxCpus='255'>pc-i440fx-jammy</machine>
			<machine maxCpus='288'>pc-q35-impish-hpb</machine>
			<machine maxCpus='288'>pc-q35-hirsute</machine>
			<machine maxCpus='288'>pc-q35-4.0.1</machine>
			<machine maxCpus='288'>pc-q35-hirsute-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-7.0</machine>
			<machine maxCpus='255' deprecated='yes'>pc-i440fx-1.6</machine>
			<machine maxCpus='288'>pc-q35-5.0</machine>
			<machine maxCpus='288'>pc-q35-2.8</machine>
			<machine maxCpus='255'>pc-i440fx-2.10</machine>
			<machine maxCpus='288'>pc-q35-3.0</machine>
			<machine maxCpus='255'>pc-i440fx-6.0</machine>
			<machine maxCpus='288'>pc-q35-zesty</machine>
			<machine maxCpus='288'>pc-q35-7.2</machine>
			<machine maxCpus='288'>pc-q35-4.0</machine>
			<machine maxCpus='288'>pc-q35-focal</machine>
			<machine maxCpus='288'>microvm</machine>
			<machine maxCpus='255'>pc-i440fx-2.3</machine>
			<machine maxCpus='255'>pc-i440fx-disco</machine>
			<machine maxCpus='255'>pc-i440fx-kinetic-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-focal-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-4.0</machine>
			<machine maxCpus='255'>pc-i440fx-groovy-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-hirsute-hpb</machine>
			<machine maxCpus='255'>pc-i440fx-5.0</machine>
			<machine maxCpus='255'>pc-i440fx-2.8</machine>
			<machine maxCpus='288'>pc-q35-6.2</machine>
			<machine maxCpus='255'>pc-i440fx-eoan</machine>
			<machine maxCpus='255'>pc-q35-2.5</machine>
			<machine maxCpus='255'>pc-i440fx-3.0</machine>
			<machine maxCpus='255'>pc-q35-yakkety</machine>
			<machine maxCpus='255'>pc-i440fx-7.2</machine>
			<machine maxCpus='288'>pc-q35-2.11</machine>
			<domain type='qemu'/>
			<domain type='kvm'/>
		</arch>
		<features>
			<acpi default='on' toggle='yes'/>
			<apic default='on' toggle='no'/>
			<cpuselection/>
			<deviceboot/>
			<disksnapshot default='on' toggle='no'/>
		</features>
		</guest>
	
	</capabilities>
	`

	var caps libvirtxml.Caps
	err := caps.Unmarshal(capsXML)
	require.Nil(t, err)

	media, err2 := installermedia.ForPath("dummy.iso")
	require.Nil(t, err2)

	actual, err3 := vmconfigurator.CreateDomainConfig(media, "whatever", caps)
	require.Nil(t, err3)

	actualXML, err4 := actual.Marshal()
	require.Nil(t, err4)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<domain type="kvm">
			<name>win11</name>
			<uuid>e7b61ffa-27f9-46c3-b1cd-67c71e469dae</uuid>
			<title>Microsoft Windows 11</title>
			<metadata>
			<boxes:gnome-boxes xmlns:boxes="https://wiki.gnome.org/Apps/Boxes">
				<os-state>installation</os-state>
				<media-id>http://microsoft.com/win/11:0</media-id>
				<media>/home/tiger/QEMU/Windows10/Windows10.iso</media>
			</boxes:gnome-boxes>
			<libosinfo:libosinfo xmlns:libosinfo="http://libosinfo.org/xmlns/libvirt/domain/1.0">
				<libosinfo:os id="http://microsoft.com/win/11"/>
			</libosinfo:libosinfo>
			</metadata>
			<memory unit="KiB">4194304</memory>
			<currentMemory unit="KiB">4194304</currentMemory>
			<vcpu placement="static">16</vcpu>
			<os>
				<type arch="x86_64" machine="pc-q35-8.0">hvm</type>
				<boot dev="cdrom"/>
				<boot dev="hd"/>
				<bootmenu enable="yes"/>
			</os>
			<features>
				<acpi/>
				<apic/>
			</features>
			<cpu mode="host-passthrough" check="none" migratable="on">
				<topology sockets="1" dies="1" cores="16" threads="1"/>
			</cpu>
			<clock offset="localtime">
				<timer name="rtc" tickpolicy="catchup"/>
				<timer name="pit" tickpolicy="delay"/>
				<timer name="hpet" present="no"/>
			</clock>
			<on_poweroff>destroy</on_poweroff>
			<on_reboot>destroy</on_reboot>
			<on_crash>destroy</on_crash>
			<pm>
				<suspend-to-mem enabled="no"/>
				<suspend-to-disk enabled="no"/>
			</pm>
			<devices>
			<emulator>/usr/bin/qemu-system-x86_64</emulator>
			<disk type="file" device="disk">
				<driver name="qemu" type="qcow2" cache="writeback"/>
				<source file="/home/tiger/.local/share/gnome-boxes/images/win11"/>
				<target dev="sda" bus="sata"/>
				<address type="drive" controller="0" bus="0" target="0" unit="0"/>
			</disk>
			<disk type="file" device="cdrom">
				<driver name="qemu" type="raw"/>
				<source file="/home/tiger/QEMU/Windows10/Windows10.iso" startupPolicy="mandatory"/>
				<target dev="hdc" bus="sata"/>
				<readonly/>
				<address type="drive" controller="0" bus="0" target="0" unit="2"/>
			</disk>
			<controller type="usb" index="0" model="qemu-xhci" ports="15">
				<address type="pci" domain="0x0000" bus="0x03" slot="0x00" function="0x0"/>
			</controller>
			<controller type="sata" index="0">
				<address type="pci" domain="0x0000" bus="0x00" slot="0x1f" function="0x2"/>
			</controller>
			<controller type="pci" index="0" model="pcie-root"/>
			<controller type="pci" index="1" model="pcie-root-port">
				<model name="pcie-root-port"/>
				<target chassis="1" port="0x10"/>
				<address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x0" multifunction="on"/>
			</controller>
			<controller type="pci" index="2" model="pcie-to-pci-bridge">
				<model name="pcie-pci-bridge"/>
				<address type="pci" domain="0x0000" bus="0x01" slot="0x00" function="0x0"/>
			</controller>
			<controller type="pci" index="3" model="pcie-root-port">
				<model name="pcie-root-port"/>
				<target chassis="3" port="0x11"/>
				<address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x1"/>
			</controller>
			<controller type="pci" index="4" model="pcie-root-port">
				<model name="pcie-root-port"/>
				<target chassis="4" port="0x12"/>
				<address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x2"/>
			</controller>
			<controller type="pci" index="5" model="pcie-root-port">
				<model name="pcie-root-port"/>
				<target chassis="5" port="0x13"/>
				<address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x3"/>
			</controller>
			<controller type="pci" index="6" model="pcie-root-port">
				<model name="pcie-root-port"/>
				<target chassis="6" port="0x14"/>
				<address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x4"/>
			</controller>
			<controller type="virtio-serial" index="0">
				<address type="pci" domain="0x0000" bus="0x04" slot="0x00" function="0x0"/>
			</controller>
			<controller type="ccid" index="0">
				<address type="usb" bus="0" port="1"/>
			</controller>
			<interface type="user">
				<mac address="52:54:00:f3:0e:57"/>
				<model type="rtl8139"/>
				<address type="pci" domain="0x0000" bus="0x02" slot="0x01" function="0x0"/>
			</interface>
			<smartcard mode="passthrough" type="spicevmc">
				<address type="ccid" controller="0" slot="0"/>
			</smartcard>
			<serial type="pty">
				<target type="isa-serial" port="0">
				<model name="isa-serial"/>
				</target>
			</serial>
			<console type="pty">
				<target type="serial" port="0"/>
			</console>
			<channel type="spicevmc">
				<target type="virtio" name="com.redhat.spice.0"/>
				<address type="virtio-serial" controller="0" bus="0" port="1"/>
			</channel>
			<channel type="spiceport">
				<source channel="org.spice-space.webdav.0"/>
				<target type="virtio" name="org.spice-space.webdav.0"/>
				<address type="virtio-serial" controller="0" bus="0" port="2"/>
			</channel>
			<input type="tablet" bus="usb">
				<address type="usb" bus="0" port="2"/>
			</input>
			<input type="mouse" bus="ps2"/>
			<input type="keyboard" bus="ps2"/>
			<graphics type="spice">
				<listen type="none"/>
				<image compression="off"/>
				<gl enable="no"/>
			</graphics>
			<sound model="ich9">
				<address type="pci" domain="0x0000" bus="0x00" slot="0x1b" function="0x0"/>
			</sound>
			<audio id="1" type="spice"/>
			<video>
				<model type="qxl" ram="65536" vram="65536" vgamem="16384" heads="1" primary="yes">
				<acceleration accel3d="no"/>
				</model>
				<address type="pci" domain="0x0000" bus="0x00" slot="0x01" function="0x0"/>
			</video>
			<redirdev bus="usb" type="spicevmc">
				<address type="usb" bus="0" port="3"/>
			</redirdev>
			<redirdev bus="usb" type="spicevmc">
				<address type="usb" bus="0" port="4"/>
			</redirdev>
			<redirdev bus="usb" type="spicevmc">
				<address type="usb" bus="0" port="5"/>
			</redirdev>
			<redirdev bus="usb" type="spicevmc">
				<address type="usb" bus="0" port="6"/>
			</redirdev>
			<watchdog model="itco" action="reset"/>
			<memballoon model="virtio">
				<address type="pci" domain="0x0000" bus="0x05" slot="0x00" function="0x0"/>
			</memballoon>
			</devices>
		</domain>
	`

	var expected libvirtxml.Domain
	err = expected.Unmarshal(expectedXML)
	require.Nil(t, err)

	require.Equal(t, expected.Type, actual.Type)
	// require.Equal(t, expected.Name, actual.Name)
	// require.Equal(t, expected.Memory, actual.Memory)
	require.Equal(t, expected.VCPU.Value, actual.VCPU.Value)

	require.Equal(t, expected.OS.Type.Arch, actual.OS.Type.Arch)
	// require.Equal(t, expected.OS.Type.Machine, actual.OS.Type.Machine)
	require.Equal(t, expected.OS.Type.Type, actual.OS.Type.Type)

	require.Equal(t, expected.Features, actual.Features)

	require.Equal(t, expected.CPU.Mode, actual.CPU.Mode)
	require.Equal(t, expected.CPU.Topology, actual.CPU.Topology)

	require.Equal(t, expected.Clock, actual.Clock)

	require.Equal(t, expected.OnPoweroff, actual.OnPoweroff)
	require.Equal(t, expected.OnReboot, actual.OnReboot)
	require.Equal(t, expected.OnCrash, actual.OnCrash)

	require.Equal(t, expected.PM, actual.PM)
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

func TestAddUSBSupport(t *testing.T) {
	domain := libvirtxml.Domain{}

	vmconfigurator.AddUSBSupport(&domain)

	require.NotNil(t, domain.Devices)
	require.NotNil(t, domain.Devices.RedirDevs)
	require.Len(t, domain.Devices.RedirDevs, 4)
	require.NotNil(t, domain.Devices.Controllers)
	require.Len(t, domain.Devices.Controllers, 1)

	actual := domain.Devices
	require.NotNil(t, actual)

	actualXML, err := xml.Marshal(domain.Devices)
	require.Nil(t, err)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<devices>
			<redirdev bus="usb" type="spicevmc">
				<address type="usb" bus="0" port="3"/>
			</redirdev>
			<redirdev bus="usb" type="spicevmc">
				<address type="usb" bus="0" port="4"/>
			</redirdev>
			<redirdev bus="usb" type="spicevmc">
				<address type="usb" bus="0" port="5"/>
			</redirdev>
			<redirdev bus="usb" type="spicevmc">
				<address type="usb" bus="0" port="6"/>
			</redirdev>
			<controller type="usb" index="0" model="qemu-xhci" ports="15">
				<address type="pci" domain="0x0000" bus="0x03" slot="0x00" function="0x0"/>
			</controller>
		</devices>
	`

	var expected libvirtxml.DomainDeviceList
	err = xml.Unmarshal([]byte(expectedXML), &expected)
	require.Nil(t, err)

	for i, actualRedirDev := range actual.RedirDevs {
		require.Equal(t, expected.RedirDevs[i].Bus, actualRedirDev.Bus)
		require.Equal(t, expected.RedirDevs[i].Source, actualRedirDev.Source)
	}

	require.Equal(t, expected.Controllers[0].Type, actual.Controllers[0].Type)
	require.Equal(t, expected.Controllers[0].Index, actual.Controllers[0].Index)
	require.Equal(t, expected.Controllers[0].Model, actual.Controllers[0].Model)
	require.Equal(t, expected.Controllers[0].USB.Port, actual.Controllers[0].USB.Port)
}

func TestCreateNetworkInterface(t *testing.T) {
	domain := libvirtxml.Domain{}

	actual := vmconfigurator.CreateNetworkInterface(&domain, false, false)

	actualXML, err := actual.Marshal()
	require.Nil(t, err)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<interface type="user">
			<mac address="52:54:00:f3:0e:57"/>
			<model type="rtl8139"/>
			<address type="pci" domain="0x0000" bus="0x02" slot="0x01" function="0x0"/>
		</interface>
	`

	var expected libvirtxml.DomainInterface
	err = expected.Unmarshal(expectedXML)
	require.Nil(t, err)

	require.NotNil(t, actual.Source)
	require.NotNil(t, actual.Source.User)

	require.Equal(t, expected.Source, actual.Source)
}

func TestCreateSpiceAgentChannel(t *testing.T) {
	actual := vmconfigurator.CreateSpiceAgentChannel()

	actualXML, err := actual.Marshal()
	require.Nil(t, err)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<channel type="spicevmc">
			<target type="virtio" name="com.redhat.spice.0"/>
			<address type="virtio-serial" controller="0" bus="0" port="1"/>
		</channel>
	`

	var expected libvirtxml.DomainChannel
	err = expected.Unmarshal(expectedXML)
	require.Nil(t, err)

	require.Equal(t, expected.Source, actual.Source)
	require.Equal(t, expected.Target, actual.Target)
}

func TestCreateGraphicDevice(t *testing.T) {
	actual := libvirtxml.DomainGraphic{
		Spice: vmconfigurator.CreateGraphicDevice(nil),
	}

	actualXML, err := actual.Marshal()
	require.Nil(t, err)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<graphics type="spice">
			<listen type="none"/>
			<image compression="off"/>
			<gl enable="no"/>
		</graphics>
	`

	var expected libvirtxml.DomainGraphic
	err = expected.Unmarshal(expectedXML)
	require.Nil(t, err)

	require.Equal(t, expected.Spice.AutoPort, actual.Spice.AutoPort)
	require.Equal(t, expected.Spice.GL, actual.Spice.GL)
	require.Equal(t, expected.Spice.Image, actual.Spice.Image)
}

func TestCreateSpiceWebDAVChannel(t *testing.T) {
	actual := vmconfigurator.CreateSpiceWebDAVChannel()

	actualXML, err := actual.Marshal()
	require.Nil(t, err)

	t.Logf("\n%s", actualXML)

	// from real XML
	expectedXML := `
		<channel type="spiceport">
			<source channel="org.spice-space.webdav.0"/>
			<target type="virtio" name="org.spice-space.webdav.0"/>
			<address type="virtio-serial" controller="0" bus="0" port="2"/>
		</channel>  
	`

	var expected libvirtxml.DomainChannel
	err = expected.Unmarshal(expectedXML)
	require.Nil(t, err)

	require.Equal(t, expected.Source, actual.Source)
	require.Equal(t, expected.Target, actual.Target)
}
