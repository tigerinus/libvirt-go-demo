package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	config, err := GetPoolConfig()

	require.Nil(t, err)

	require.NotEmpty(t, config)
}
