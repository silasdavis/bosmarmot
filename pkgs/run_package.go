package pkgs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monax/bosmarmot/pkgs/definitions"
	"github.com/monax/bosmarmot/pkgs/jobs"
	log "github.com/sirupsen/logrus"
)

func RunPackage(do *definitions.Packages) error {
	var gotwd string
	if do.Path == "" {
		var err error
		gotwd, err = os.Getwd()
		if err != nil {
			return err
		}
		do.Path = gotwd
	}

	// note: [zr] this could be problematic with a combo of
	// other flags, however, at least the --dir flag isn't
	// completely broken now
	if do.Path != gotwd {
		originalYAMLPath := do.YAMLPath

		// if --dir given, assume *.yaml is in there
		do.YAMLPath = filepath.Join(do.Path, originalYAMLPath)

		// if do.YAMLPath does not exist, try $pwd
		if _, err := os.Stat(do.YAMLPath); os.IsNotExist(err) {
			do.YAMLPath = filepath.Join(gotwd, originalYAMLPath)
		}

		// if it still cannot be found, abort
		if _, err := os.Stat(do.YAMLPath); os.IsNotExist(err) {
			return fmt.Errorf("could not find jobs file (%s), ensure correct used of the --file flag",
				do.YAMLPath)
		}

		if do.BinPath == "./bin" {
			do.BinPath = filepath.Join(do.Path, "bin")
		}
		if do.ABIPath == "./abi" {
			do.ABIPath = filepath.Join(do.Path, "abi")
		}
		// TODO enable this feature
		// if do.ContractsPath == "./contracts" {
		//do.ContractsPath = filepath.Join(do.Path, "contracts")
		//}
	}

	// useful for debugging
	printPathPackage(do)

	var err error
	// Load the package if it doesn't exist
	if do.Package == nil {
		do.Package, err = LoadPackage(do.YAMLPath)
		if err != nil {
			return err
		}
	}

	if do.Path != gotwd {
		for _, job := range do.Package.Jobs {
			if job.Deploy != nil {
				job.Deploy.Contract = filepath.Join(do.Path, job.Deploy.Contract)
			}
		}
	}

	return jobs.RunJobs(do)
}

func printPathPackage(do *definitions.Packages) {
	log.WithField("=>", do.ChainURL).Info("With ChainURL")
	log.WithField("=>", do.Signer).Info("Using Signer at")
}
