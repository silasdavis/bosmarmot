package commands

import (
	"fmt"

	"github.com/monax/bosmarmot/monax/pkgs"
	"github.com/monax/bosmarmot/monax/util"

	"github.com/monax/bosmarmot/monax/keys"
	"github.com/spf13/cobra"
)

var Packages = &cobra.Command{
	Use:   "pkgs",
	Short: "deploy, test, and manage your smart contract packages",
	Long: `the pkgs subcommand is used to test and deploy
smart contract packages for use by your application`,
	Run: func(cmd *cobra.Command, args []string) {
		util.IfExit(ArgCheck(0, "eq", cmd, args))
		if do.DefaultAddr == "" { // note that this is not strictly necessary since the addr can be set in the epm.yaml.
			util.IfExit(fmt.Errorf("please provide the address to deploy from with --address"))
		}

		util.IfExit(pkgs.RunPackage(do))
	},
}

func buildPackagesCommand() {
	addPackagesFlags()
}

func addPackagesFlags() {
	Packages.Flags().StringVarP(&do.ChainURL, "chain-url", "", "tcp://localhost:46657", "chain-url to be used in tcp://IP:PORT format (only necessary for cluster and remote operations)")
	Packages.Flags().StringVarP(&do.Signer, "keys", "s", defaultSigner(), "IP:PORT of keys daemon which jobs should use")
	Packages.Flags().StringVarP(&do.Path, "dir", "i", "", "root directory of app (will use $pwd by default)")
	Packages.Flags().StringVarP(&do.DefaultOutput, "output", "o", "epm.output.json", "filename for jobs output file. by default, this name will reflect the name passed in on the optional [--file]")
	Packages.Flags().StringVarP(&do.YAMLPath, "file", "f", "epm.yaml", "path to package file which jobs should use. if also using the --dir flag, give the relative path to jobs file, which should be in the same directory")
	Packages.Flags().StringSliceVarP(&do.DefaultSets, "set", "e", []string{}, "default sets to use; operates the same way as the [set] jobs, only before the jobs file is ran (and after default address")
	// the package manager does not use this flag!
	// Packages.Flags().StringVarP(&do.ContractsPath, "contracts-path", "p", "./contracts", "path to the contracts jobs should use")
	Packages.Flags().StringVarP(&do.BinPath, "bin-path", "", "./bin", "path to the bin directory jobs should use when saving binaries after the compile process")
	Packages.Flags().StringVarP(&do.ABIPath, "abi-path", "", "./abi", "path to the abi directory jobs should use when saving ABIs after the compile process")
	Packages.Flags().StringVarP(&do.DefaultGas, "gas", "g", "1111111111", "default gas to use; can be overridden for any single job")
	Packages.Flags().StringVarP(&do.DefaultAddr, "address", "a", "", "default address to use; operates the same way as the [account] job, only before the epm file is ran")
	Packages.Flags().StringVarP(&do.DefaultFee, "fee", "n", "9999", "default fee to use")
	Packages.Flags().StringVarP(&do.DefaultAmount, "amount", "u", "9999", "default amount to use")
	Packages.Flags().BoolVarP(&do.Overwrite, "overwrite", "t", true, "overwrite jobs of the same name")
}

func defaultSigner() string {
	return keys.DefaultKeysURL()
}
