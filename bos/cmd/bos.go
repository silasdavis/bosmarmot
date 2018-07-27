package commands

import (
	"bytes"
	"fmt"

	"github.com/monax/bosmarmot/bos"
	"github.com/monax/bosmarmot/bos/def"
	"github.com/monax/bosmarmot/bos/util"
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
		log.SetFormatter(new(PlainFormatter))
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
var do *def.Packages

// Controls the execution sequence of the cobra global runner
func Execute() {
	do = new(def.Packages)
	AddGlobalFlags()
	AddCommands()
	util.IfExit(BosCmd.Execute())
}

// Flags that are to be used by commands are handled by the Do struct
func AddGlobalFlags() {
	BosCmd.Flags().StringVarP(&do.ChainURL, "chain-url", "u", "127.0.0.1:10997",
		"chain-url to be used in IP:PORT format")

	BosCmd.Flags().StringVarP(&do.Signer, "keys", "s", "",
		"IP:PORT of Burrow GRPC service which jobs should or otherwise transaction submitted unsigned for mempool signing in Burrow")

	BosCmd.Flags().BoolVarP(&do.MempoolSigning, "mempool-signing", "p", false,
		"Use Burrows own keys connection to sign transactions - means that Burrow instance must have access to input account keys. "+
			"Sequence numbers are set as transactions enter the mempool so concurrent transactions can be sent from same inputs.")

	BosCmd.Flags().StringVarP(&do.Path, "dir", "i", "",
		"root directory of app (will use pwd by default)")

	BosCmd.Flags().StringVarP(&do.DefaultOutput, "output", "o", "epm.output.json",
		"filename for jobs output file. by default, this name will reflect the name passed in on the optional [--file]")

	BosCmd.Flags().StringVarP(&do.YAMLPath, "file", "f", "epm.yaml",
		"path to package file which jobs should use. if also using the --dir flag, give the relative path to jobs file, which should be in the same directory")

	BosCmd.Flags().StringSliceVarP(&do.DefaultSets, "set", "e", []string{},
		"default sets to use; operates the same way as the [set] jobs, only before the jobs file is ran (and after default address")

	BosCmd.Flags().StringVarP(&do.BinPath, "bin-path", "b", "[dir]/bin",
		"path to the bin directory jobs should use when saving binaries after the compile process defaults to --dir + /bin")

	BosCmd.Flags().StringVarP(&do.DefaultGas, "gas", "g", "1111111111",
		"default gas to use; can be overridden for any single job")

	BosCmd.Flags().StringVarP(&do.Address, "address", "a", "",
		"default address to use; operates the same way as the [account] job, only before the epm file is ran")

	BosCmd.Flags().StringVarP(&do.DefaultFee, "fee", "n", "9999",
		"default fee to use")

	BosCmd.Flags().StringVarP(&do.DefaultAmount, "amount", "m", "9999",
		"default amount to use")

	BosCmd.Flags().BoolVarP(&do.Verbose, "verbose", "v", false,
		"verbose output")

	BosCmd.Flags().BoolVarP(&do.Debug, "debug", "d", false,
		"debug level output")
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

type PlainFormatter struct{}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.appendMessage(b, entry.Message)
	for _, key := range keys {
		f.appendMessageData(b, key, entry.Data[key])
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *PlainFormatter) appendMessage(b *bytes.Buffer, message string) {
	fmt.Fprintf(b, "%-44s", message)
}

func (f *PlainFormatter) appendMessageData(b *bytes.Buffer, key string, value interface{}) {
	switch key {
	case "":
		b.WriteString("=> ")
	case "=>":
		b.WriteString(key)
		b.WriteByte(' ')
	default:
		b.WriteString(key)
		b.WriteString(" => ")
	}
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	b.WriteString(stringVal)
}
