package loader

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/monax/bosmarmot/bos/def"
		"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
)

const testPackageYAML = `jobs:

- name: AddValidators
  govern-account:
    source: foo
    target: bar
    permissions: []
    roles: []

- name: nameRegTest1
  register:
    name: $val1
    data: $val2
    amount: $to_save
    fee: $MinersFee
`

func TestUnmarshal(t *testing.T) {
	pkgs := viper.New()
	pkgs.SetConfigType("yaml")
	err := pkgs.ReadConfig(bytes.NewBuffer([]byte(testPackageYAML)))
	require.NoError(t, err)
	do := new(def.Package)

	err = pkgs.UnmarshalExact(do)
	require.NoError(t, err)
	yamlOut, err := yaml.Marshal(do)
	require.NoError(t, err)
	assert.True(t, len(yamlOut) > 100, "should marshal some yaml")

	doOut := new(def.Package)
	err = yaml.Unmarshal(yamlOut, doOut)
	require.NoError(t, err)
	assert.Equal(t, do, doOut)
}
