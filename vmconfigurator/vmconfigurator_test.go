package vmconfigurator_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tigerinus/libvirt-go-demo/vmconfigurator"
)

func TestXxx(t *testing.T) {
	config, err := vmconfigurator.GetPoolConfig()

	require.Nil(t, err)

	require.NotEmpty(t, config)
}
