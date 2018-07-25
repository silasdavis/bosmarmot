package def

import (
	"errors"
	"fmt"

	"regexp"
	"strconv"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/hyperledger/burrow/crypto"
)

var Address = validation.NewStringRule(IsAddress, "must be valid 20 byte hex-encoded string")

var Placeholder = validation.Match(regexp.MustCompile(`\$[[:alnum:]]+`))

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

var Uint64 = validation.By(func(value interface{}) error {
	str, err := validation.EnsureString(value)
	if err != nil {
		return fmt.Errorf("should be number but '%v' is not a string", value)
	}
	_, err = strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fmt.Errorf("should be a 64 bit unsigned integer: ")
	}
	return nil
})

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

type BoolRule struct {
	isValid func(value interface{}) bool
	message string
}

func NewBoolRule(isValid func(value interface{}) bool, message string, args ...interface{}) *BoolRule {
	return &BoolRule{
		isValid: isValid,
		message: fmt.Sprintf(message, args...),
	}
}

func (r *BoolRule) Validate(value interface{}) error {
	if r.isValid(value) {
		return nil
	}
	return errors.New(r.message)
}

func (r *BoolRule) Error(message string, args ...interface{}) *BoolRule {
	r.message = fmt.Sprintf(message, args...)
	return r
}
