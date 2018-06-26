package project

import (
	"github.com/monax/relic"
)

// Can be used to set the commit hash version of the binary at build time with:
// `go build -ldflags "-X github.com/hyperledger/burrow/project.commit=$(git rev-parse --short HEAD)" ./cmd/burrow`

var commit = ""

func Commit() string {
	return commit
}

func FullVersion() string {
	version := History.CurrentVersion().String()
	if commit != "" {
		return version + "+commit." + commit
	}
	return version
}

// The releases described by version string and changes, newest release first.
// The current release is taken to be the first release in the slice, and its
// version determines the single authoritative version for the next release.
//
// To cut a new release add a release to the front of this slice then run the
// release tagging script: ./scripts/tag_release.sh
var History relic.ImmutableHistory = relic.NewHistory("Bosmarmot").MustDeclareReleases(
	"0.3.0",
	`Add meta job; simplify run_packages significantly; js upgrades; routine burrow compatibility upgrades`,
	"0.2.1",
	`Fix release to harmonize against burrow versions > 0.18.0`,
	"0.2.0",
	`Simplify repository by removing latent tooling and consolidating compilers and bos,
as well as removing keys completely which have been migrated to burrow`,
	"0.1.0",
	`Major release of Bosmarmot tooling including updated javascript libraries for Burrow 0.18.*`,
	"0.0.1",
	`Initial Bosmarmot combining and refactoring Monax tooling, including:
- The monax tool (just 'monax pkgs do')
- The monax-keys signing daemon
- Monax compilers
- A basic legacy-contracts.js integration test (merging in JS libs is pending)`,
)
