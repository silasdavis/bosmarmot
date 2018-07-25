# Bosmarmot Changelog
## Version 0.4.0
Big improvements to performance across bos, including:
- Implement our own ABI library rather than relying on 2.7M lines of go-ethereum code for a fairly simple library.
- Migrate bos to leveraging burrow's GPRC framework entirely.
	

## Version 0.3.0
Add meta job; simplify run_packages significantly; js upgrades; routine burrow compatibility upgrades

## Version 0.2.1
Fix release to harmonize against burrow versions > 0.18.0

## Version 0.2.0
Simplify repository by removing latent tooling and consolidating compilers and bos,
as well as removing keys completely which have been migrated to burrow

## Version 0.1.0
Major release of Bosmarmot tooling including updated javascript libraries for Burrow 0.18.*

## Version 0.0.1
Initial Bosmarmot combining and refactoring Monax tooling, including:
- The monax tool (just 'monax pkgs do')
- The monax-keys signing daemon
- Monax compilers
- A basic legacy-contracts.js integration test (merging in JS libs is pending)

