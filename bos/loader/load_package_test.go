package loader

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

const testPackageYAML = `
jobs:
- name: AddValidator
  govern:
    
  
`

func TestUnmarshal(t *testing.T) {
	viper.New()
	err := viper.ReadConfig(bytes.NewBuffer([]byte(testPackageYAML)))
	require.NoError(t, err)
}
