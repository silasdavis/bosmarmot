package def

import "github.com/go-ozzo/ozzo-validation"

type Packages struct {
	Address       string   `mapstructure:"," json:"," yaml:"," toml:","`
	BinPath       string   `mapstructure:"," json:"," yaml:"," toml:","`
	ChainURL      string   `mapstructure:"," json:"," yaml:"," toml:","`
	CurrentOutput string   `mapstructure:"," json:"," yaml:"," toml:","`
	Debug         bool     `mapstructure:"," json:"," yaml:"," toml:","`
	DefaultAmount string   `mapstructure:"," json:"," yaml:"," toml:","`
	DefaultFee    string   `mapstructure:"," json:"," yaml:"," toml:","`
	DefaultGas    string   `mapstructure:"," json:"," yaml:"," toml:","`
	DefaultOutput string   `mapstructure:"," json:"," yaml:"," toml:","`
	DefaultSets   []string `mapstructure:"," json:"," yaml:"," toml:","`
	Path          string   `mapstructure:"," json:"," yaml:"," toml:","`
	Signer        string   `mapstructure:"," json:"," yaml:"," toml:","`
	Verbose       bool     `mapstructure:"," json:"," yaml:"," toml:","`
	YAMLPath      string   `mapstructure:"," json:"," yaml:"," toml:","`

	Package *Package
	Client
}

func (do *Packages) Validate() error {
	return validation.ValidateStruct(do,
		validation.Field(&do.Address, Address),
	)
}

func (do *Packages) Dial() error {
	return do.Client.Dial(do.ChainURL, do.Signer)
}

type Package struct {
	// from epm
	Account string
	Jobs    []*Job
}

func (pkg *Package) Validate() error {
	return validation.ValidateStruct(pkg,
		validation.Field(&pkg.Account, Address),
		validation.Field(&pkg.Jobs),
	)
}
