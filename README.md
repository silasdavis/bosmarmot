# Bosmarmot

|[![GoDoc](https://godoc.org/github.com/bosmarmot?status.png)](https://godoc.org/github.com/monax/bosmarmot/bos/cmd) | Linux |
|---|-------|
| Master | [![Circle CI](https://circleci.com/gh/monax/bosmarmot/tree/master.svg?style=svg)](https://circleci.com/gh/monax/bosmarmot/tree/master) |
| Develop | [![Circle CI (develop)](https://circleci.com/gh/monax/bosmarmot/tree/develop.svg?style=svg)](https://circleci.com/gh/monax/bosmarmot/tree/develop) |


Monax maintains Bosmarmot a satellite monorepo to 
[Burrow](https://github.com/hyperledger/burrow) that contains:

* burrow.js - a Javascript client library for interacting with Burrow smart contracts
* Vent - an EVM event to SQL database mapping layer - currently under development.

If you are looking for `bos`, this has been replaced by `burrow deploy` in
[Burrow](https://github.com/hyperledger/burrow).

## Working with Javascript

Currently the javascript libraries are being rebuilt. The master branch of this repository works against the master branch of burrow.

Please use the versions within the package.json of this repo on master branch for fully compatible and tested versions.
