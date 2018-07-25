package def

import (
	"regexp"

	"reflect"

	"fmt"

	"github.com/go-ozzo/ozzo-validation"
)

//TODO: Interface all the jobs, determine if they should remain in definitions or get their own package

type Job struct {
	// Name of the job
	Name string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	// Not marshalled
	Result string
	// For multiple values
	Variables []*Variable
	// Sets/Resets the primary account to use
	Account *Account `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// Set an arbitrary value
	Set *SetJob `mapstructure:"set" json:"set" yaml:"set" toml:"set"`
	// Run a sequence of other epm.yamls
	Meta *Meta `mapstructure:"meta" json:"meta" yaml:"meta" toml:"meta"`
	// Issue a governance transaction
	Govern *Govern `mapstructure:"govern" json:"govern" yaml:"govern" toml:"govern"`
	// Contract compile and send to the chain functions
	Deploy *Deploy `mapstructure:"deploy" json:"deploy" yaml:"deploy" toml:"deploy"`
	// Send tokens from one account to another
	Send *Send `mapstructure:"send" json:"send" yaml:"send" toml:"send"`
	// Utilize monax:db's native name registry to register a name
	RegisterName *RegisterName `mapstructure:"register" json:"register" yaml:"register" toml:"register"`
	// Sends a transaction which will update the permissions of an account. Must be sent from an account which
	// has root permissions on the blockchain (as set by either the genesis.json or in a subsequence transaction)
	Permission *Permission `mapstructure:"permission" json:"permission" yaml:"permission" toml:"permission"`
	// Sends a transaction to a contract. Will utilize monax-abi under the hood to perform all of the heavy lifting
	Call *Call `mapstructure:"call" json:"call" yaml:"call" toml:"call"`
	// Wrapper for mintdump dump. WIP
	DumpState *DumpState `mapstructure:"dump-state" json:"dump-state" yaml:"dump-state" toml:"dump-state"`
	// Wrapper for mintdum restore. WIP
	RestoreState *RestoreState `mapstructure:"restore-state" json:"restore-state" yaml:"restore-state" toml:"restore-state"`
	// Sends a "simulated call" to a contract. Predominantly used for accessor functions ("Getters" within contracts)
	QueryContract *QueryContract `mapstructure:"query-contract" json:"query-contract" yaml:"query-contract" toml:"query-contract"`
	// Queries information from an account.
	QueryAccount *QueryAccount `mapstructure:"query-account" json:"query-account" yaml:"query-account" toml:"query-account"`
	// Queries information about a name registered with monax:db's native name registry
	QueryName *QueryName `mapstructure:"query-name" json:"query-name" yaml:"query-name" toml:"query-name"`
	// Queries information about the validator set
	QueryVals *QueryVals `mapstructure:"query-vals" json:"query-vals" yaml:"query-vals" toml:"query-vals"`
	// Makes and assertion (useful for testing purposes)
	Assert *Assert `mapstructure:"assert" json:"assert" yaml:"assert" toml:"assert"`
}

func (job *Job) Validate() error {
	payload, err := job.Validatable()
	if err != nil {
		return err
	}
	return validation.ValidateStruct(job,
		validation.Field(&job.Name, validation.Required, validation.Match(regexp.MustCompile("[[:word:]]+")).
			Error("must contain word characters; alphanumeric plus underscores/hyphens")),
		validation.Field(&job.Result, NewBoolRule(IsOmitted, "internally reserved and should be removed")),
		validation.Field(&job.Variables, NewBoolRule(IsOmitted, "internally reserved and should be removed")),
		validation.Field(payload),
	)
}

var validatableType = reflect.TypeOf((*validation.Validatable)(nil)).Elem()

// Ensures only one Job payload is set and returns a pointer to that field or an error if none or multiple
// job payload fields are set
func (job *Job) Validatable() (interface{}, error) {
	rv := reflect.ValueOf(job).Elem()
	rt := rv.Type()

	payloadIndex := -1
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Type.Implements(validatableType) && !rv.Field(i).IsNil() {
			if payloadIndex >= 0 {
				return nil, fmt.Errorf("only one Job payload field should be set, but both '%v' and '%v' are set",
					rt.Field(payloadIndex).Name, rt.Field(i).Name)
			}
			payloadIndex = i
		}
	}
	if payloadIndex == -1 {
		return nil, fmt.Errorf("Job has no payload, please set at least one job value")
	}

	return rv.Field(payloadIndex).Addr().Interface(), nil
}
