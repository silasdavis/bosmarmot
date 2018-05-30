package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/monax/bosmarmot/pkgs/definitions"

	acm "github.com/hyperledger/burrow/account"
	"github.com/hyperledger/burrow/client"
	"github.com/hyperledger/burrow/logging"
)

func GetBlockHeight(do *definitions.Packages) (latestBlockHeight uint64, err error) {
	nodeClient := client.NewBurrowNodeClient(do.ChainURL, logging.NewNoopLogger())
	// NOTE: NodeInfo is no longer exposed through Status();
	// other values are currently not use by the package manager
	_, _, _, latestBlockHeight, _, err = nodeClient.Status()
	if err != nil {
		return 0, err
	}
	// set return values
	return
}

func AccountsInfo(account, field string, do *definitions.Packages) (string, error) {

	address, err := acm.AddressFromHexString(account)
	if err != nil {
		return "", fmt.Errorf("Account Addr %s is improper hex: %v", account, err)
	}
	nodeClient := client.NewBurrowNodeClient(do.ChainURL, logging.NewNoopLogger())

	r, err := nodeClient.GetAccount(address)
	if err != nil {
		return "", err
	}
	if r == nil {
		return "", fmt.Errorf("Account %s does not exist", account)
	}

	var s string
	if strings.Contains(field, "permissions") {
		// TODO: [ben] resolve conflict between explicit types and json better

		fields := strings.Split(field, ".")

		if len(fields) > 1 {
			switch fields[1] {
			case "roles":
				s = strings.Join(r.Permissions().Roles, ",")
			case "base", "perms":
				s = strconv.Itoa(int(r.Permissions().Base.Perms))
			case "set":
				s = strconv.Itoa(int(r.Permissions().Base.SetBit))
			}
		}
	} else if field == "balance" {
		s = itoaU64(r.Balance())
	}

	if err != nil {
		return "", err
	}

	return s, nil
}

func NamesInfo(name, field string, do *definitions.Packages) (string, error) {
	nodeClient := client.NewBurrowNodeClient(do.ChainURL, logging.NewNoopLogger())
	owner, data, expirationBlock, err := nodeClient.GetName(name)
	if err != nil {
		return "", err
	}

	switch strings.ToLower(field) {
	case "name":
		return name, nil
	case "owner":
		return owner.String(), nil
	case "data":
		return data, nil
	case "expires":
		return itoaU64(expirationBlock), nil
	default:
		return "", fmt.Errorf("Field %s not recognized", field)
	}
}

func ValidatorsInfo(field string, do *definitions.Packages) (string, error) {
	nodeClient := client.NewBurrowNodeClient(do.ChainURL, logging.NewNoopLogger())
	_, bondedValidators, unbondingValidators, err := nodeClient.ListValidators()
	if err != nil {
		return "", err
	}

	vals := []string{}
	switch strings.ToLower(field) {
	case "bonded_validators":
		for _, v := range bondedValidators {
			vals = append(vals, v.Address().String())
		}
	case "unbonding_validators":
		for _, v := range unbondingValidators {
			vals = append(vals, v.Address().String())
		}
	default:
		return "", fmt.Errorf("Field %s not recognized", field)
	}
	return strings.Join(vals, ","), nil
}

func itoaU64(i uint64) string {
	return strconv.FormatUint(i, 10)
}
