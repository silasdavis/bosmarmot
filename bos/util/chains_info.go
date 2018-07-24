package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/monax/bosmarmot/bos/def"

	"github.com/hyperledger/burrow/crypto"
)

func GetBlockHeight(do *def.Packages) (latestBlockHeight uint64, err error) {
	stat, err := do.Status()
	if err != nil {
		return 0, err
	}
	return stat.LatestBlockHeight, nil
}

func AccountsInfo(account, field string, do *def.Packages) (string, error) {
	address, err := crypto.AddressFromHexString(account)
	if err != nil {
		return "", err
	}
	acc, err := do.GetAccount(address)
	if err != nil {
		return "", err
	}
	if acc == nil {
		return "", fmt.Errorf("Account %s does not exist", account)
	}

	var s string
	if strings.Contains(field, "permissions") {
		fields := strings.Split(field, ".")

		if len(fields) > 1 {
			switch fields[1] {
			case "roles":
				s = strings.Join(acc.Permissions.Roles, ",")
			case "base", "perms":
				s = strconv.Itoa(int(acc.Permissions.Base.Perms))
			case "set":
				s = strconv.Itoa(int(acc.Permissions.Base.SetBit))
			}
		}
	} else if field == "balance" {
		s = itoaU64(acc.Balance)
	}

	if err != nil {
		return "", err
	}

	return s, nil
}

func NamesInfo(name, field string, do *def.Packages) (string, error) {
	entry, err := do.GetName(name)
	if err != nil {
		return "", err
	}

	switch strings.ToLower(field) {
	case "name":
		return name, nil
	case "owner":
		return entry.Owner.String(), nil
	case "data":
		return entry.Data, nil
	case "expires":
		return itoaU64(entry.Expires), nil
	default:
		return "", fmt.Errorf("Field %s not recognized", field)
	}
}

func ValidatorsInfo(field string, do *def.Packages) (string, error) {
	// Currently there is no notion of 'unbonding validators' we can revisit what should go here or whether this deserves
	// to exist as a job
	if field == "bonded_validators" {
		set, err := do.GetValidatorSet()
		if err != nil {
			return "", err
		}
		return set.String(), nil
	}
	return "", nil
}

func itoaU64(i uint64) string {
	return strconv.FormatUint(i, 10)
}
