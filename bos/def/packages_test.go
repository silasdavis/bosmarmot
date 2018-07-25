package def

import (
	"testing"

	"github.com/hyperledger/burrow/crypto"
	"github.com/stretchr/testify/require"
)

func TestPackage_Validate(t *testing.T) {
	pkgs := &Package{
		Jobs: []*Job{{
			Name: "CallJob",
			Call: &Call{
				Sequence: "13",
			},
		}},
	}
	err := pkgs.Validate()
	require.NoError(t, err)

	address := crypto.Address{3, 4}.String()
	pkgs.Jobs = append(pkgs.Jobs, &Job{
		Name: "Foo",
		Account: &Account{
			Address: address,
		},
	})
	err = pkgs.Validate()
	require.NoError(t, err)

	// cannot set two job fields
	pkgs.Jobs[1].QueryAccount = &QueryAccount{
		Account: address,
		Field:   "Foo",
	}
	err = pkgs.Validate()
	require.Error(t, err)
}
