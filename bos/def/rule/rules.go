package rule

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/permission"
)

var VariableRegex = regexp.MustCompile(`(^|\s|:)\$([a-zA-Z0-9_.]+)`)

var exampleAddress = acm.GeneratePrivateAccountFromSecret("marmot")

// Rules
var (
	Placeholder = validation.Match(VariableRegex)

	Address = validation.NewStringRule(IsAddress,
		fmt.Sprintf("must be valid 20 byte hex-encoded string like '%v'", exampleAddress))

	AddressOrPlaceholder = Or(Placeholder, Address)

	Relation = validation.In("eq", "ne", "ge", "gt", "le", "lt", "==", "!=", ">=", ">", "<=", "<")

	HexOrPlaceholder = Or(Placeholder, is.Hexadecimal)

	PermissionOrPlaceholder = Or(Placeholder, Permission)

	Permission = validation.By(func(value interface{}) error {
		str, err := validation.EnsureString(value)
		if err != nil {
			return fmt.Errorf("should be a permission name")
		}
		_, err = permission.PermStringToFlag(str)
		if err != nil {
			return err
		}
		return nil
	})

	Uint64OrPlaceholder = Or(Placeholder, Uint64)

	Uint64 = validation.By(func(value interface{}) error {
		str, err := validation.EnsureString(value)
		if err != nil {
			return fmt.Errorf("should be a numeric string but '%v' is not a string", value)
		}
		_, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			return fmt.Errorf("should be a 64 bit unsigned integer: ")
		}
		return nil
	})
)

type orRule struct {
	rules []validation.Rule
}

func (orr *orRule) Validate(value interface{}) error {
	errs := make([]error, len(orr.rules))
	for i, r := range orr.rules {
		errs[i] = r.Validate(value)
		if errs[i] == nil {
			return nil
		}
	}
	return fmt.Errorf("did not validate any requirements: %v", errs)
}

func Or(rules ...validation.Rule) *orRule {
	return &orRule{
		rules: rules,
	}
}

func IsAddress(value string) bool {
	_, err := crypto.AddressFromHexString(value)
	return err == nil
}

// Returns true IFF value is zero value or has length 0
func IsOmitted(value interface{}) bool {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return true
	}
	// Accept and empty slice or map
	length, err := validation.LengthOfValue(value)
	if err == nil && length == 0 {
		return true
	}
	return false
}

type boolRule struct {
	isValid func(value interface{}) bool
	message string
}

func New(isValid func(value interface{}) bool, message string, args ...interface{}) *boolRule {
	return &boolRule{
		isValid: isValid,
		message: fmt.Sprintf(message, args...),
	}
}

func (r *boolRule) Validate(value interface{}) error {
	if r.isValid(value) {
		return nil
	}
	return errors.New(r.message)
}

func (r *boolRule) Error(message string, args ...interface{}) *boolRule {
	r.message = fmt.Sprintf(message, args...)
	return r
}
