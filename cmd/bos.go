package commands

import (
	"fmt"

	"github.com/monax/bosmarmot/pkgs"
	"github.com/monax/bosmarmot/pkgs/definitions"
	"github.com/monax/bosmarmot/pkgs/util"
	"github.com/monax/bosmarmot/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Defines the root command
var BosCmd = &cobra.Command{
	Use: "bos COMMAND [FLAG ...]",
	Long: `bos is an application for deploying and testing packages to Hyperledger Burrow.

It is used to test and deploy smart contract packages for use by your application.

Bos will perform the required functionality included in a package definition file.

Made with <3 by Monax Industries.
` + "\nVersion:\n  " + project.History.CurrentVersion().String(),
	Run: func(cmd *cobra.Command, args []string) {
		util.IfExit(ArgCheck(0, "eq", cmd, args))
		log.SetLevel(log.WarnLevel)
		if do.Verbose {
			log.SetLevel(log.InfoLevel)
		} else if do.Debug {
			log.SetLevel(log.DebugLevel)
		}
		util.IfExit(pkgs.RunPackage(do))
	},
}

// Global Do struct
var do *definitions.Packages

// Controls the execution sequence of the cobra global runner
func Execute() {
	do = definitions.NewPackage()
	AddGlobalFlags()
	AddCommands()
	util.IfExit(BosCmd.Execute())
}

// Flags that are to be used by commands are handled by the Do struct
func AddGlobalFlags() {
	BosCmd.Flags().StringVarP(&do.ChainURL, "chain-url", "u", "tcp://localhost:46657", "chain-url to be used in tcp://IP:PORT format (only necessary for cluster and remote operations)")
	BosCmd.Flags().StringVarP(&do.Signer, "keys", "s", "tcp://localhost:46657", "tcp://IP:PORT of keys daemon which jobs should use")
	BosCmd.Flags().StringVarP(&do.Path, "dir", "i", "", "root directory of app (will use pwd by default)")
	BosCmd.Flags().StringVarP(&do.DefaultOutput, "output", "o", "epm.output.json", "filename for jobs output file. by default, this name will reflect the name passed in on the optional [--file]")
	BosCmd.Flags().StringVarP(&do.YAMLPath, "file", "f", "epm.yaml", "path to package file which jobs should use. if also using the --dir flag, give the relative path to jobs file, which should be in the same directory")
	BosCmd.Flags().StringSliceVarP(&do.DefaultSets, "set", "e", []string{}, "default sets to use; operates the same way as the [set] jobs, only before the jobs file is ran (and after default address")
	BosCmd.Flags().StringVarP(&do.BinPath, "bin-path", "b", "./bin", "path to the bin directory jobs should use when saving binaries after the compile process")
	BosCmd.Flags().StringVarP(&do.ABIPath, "abi-path", "p", "./abi", "path to the abi directory jobs should use when saving ABIs after the compile process")
	BosCmd.Flags().StringVarP(&do.DefaultGas, "gas", "g", "1111111111", "default gas to use; can be overridden for any single job")
	BosCmd.Flags().StringVarP(&do.Address, "address", "a", "", "default address to use; operates the same way as the [account] job, only before the epm file is ran")
	BosCmd.Flags().StringVarP(&do.DefaultFee, "fee", "n", "9999", "default fee to use")
	BosCmd.Flags().StringVarP(&do.DefaultAmount, "amount", "m", "9999", "default amount to use")
	BosCmd.Flags().BoolVarP(&do.Verbose, "verbose", "v", false, "verbose output")
	BosCmd.Flags().BoolVarP(&do.Debug, "debug", "d", false, "debug level output")
}

// Define the sub-commands
func AddCommands() {
	BosCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print Version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(project.FullVersion())
		},
	})
	BosCmd.SetHelpCommand(Help)
	BosCmd.SetHelpTemplate(helpTemplate)
}

// Utility helpers
func ArgCheck(num int, comp string, cmd *cobra.Command, args []string) error {
	switch comp {
	case "eq":
		if len(args) != num {
			cmd.Help()
			return fmt.Errorf("\n**Note** you sent our marmots the wrong number of arguments.\nPlease send the marmots %d arguments only.", num)
		}
	case "ge":
		if len(args) < num {
			cmd.Help()
			return fmt.Errorf("\n**Note** you sent our marmots the wrong number of arguments.\nPlease send the marmots at least %d argument(s).", num)
		}
	}
	return nil
}
