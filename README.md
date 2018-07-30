# Bosmarmot

|[![GoDoc](https://godoc.org/github.com/bosmarmot?status.png)](https://godoc.org/github.com/monax/bosmarmot/bos/cmd) | Linux |
|---|-------|
| Master | [![Circle CI](https://circleci.com/gh/monax/bosmarmot/tree/master.svg?style=svg)](https://circleci.com/gh/monax/bosmarmot/tree/master) |
| Develop | [![Circle CI (develop)](https://circleci.com/gh/monax/bosmarmot/tree/develop.svg?style=svg)](https://circleci.com/gh/monax/bosmarmot/tree/develop) |

Bosmarmot is a monorepo containing condensed and updated versions of the basic tooling required to interact with a [Burrow](https://github.com/hyperledger/burrow) chain.

It also contains the interpreter for the burrow packages specification language (previously known as 'epm'). This README will cover setting up a Burrow chain with the bosmarmot tooling from start to finish.

## Install

We're going to need three (3) binaries:

```
burrow
bos
solc
```

First, ensure you have `go` installed and `$GOPATH` set

For `burrow`:

```
go get github.com/hyperledger/burrow
cd $GOPATH/src/github.com/hyperledger/burrow
make build
```

which will put the `burrow` binary in `/bin`. Move it onto your `$PATH`

For `bos`:

```
go get github.com/monax/bosmarmot
cd $GOPATH/src/github.com/monax/bosmarmot
make build
```

and move these onto your `$PATH` as well.

To install the solidity compiler - `solc` - see [here](https://solidity.readthedocs.io/en/develop/installing-solidity.html) for platform specific instructions.

## Configure

The end result will be a `burrow.toml` which contains the genesis spec and burrow configuration options required when starting the `burrow` node.

### Accounts

First, let's create some accounts. In this case, we're creating one of each a Participant and Full account:

```
burrow spec --participant-accounts=1 --full-accounts=1 > genesis-spec.json
```

and writing the output to an `genesis-spec.json`. This file should look like:

```
{
	"Accounts": [
		{
			"Name": "Full_0",
			"Amounts": [
				{
					"Type": "Native",
					"Amount": 99999999999999
				},
				{
					"Type": "Power",
					"Amount": 9999999999
				}
			],
			"Permissions": [
				"all"
			]
		},
		{
			"Name": "Participant_0",
			"Amounts": [
				{
					"Type": "Native",
					"Amount": 9999999999
				}
			],
			"Permissions": [
				"send",
				"call",
				"name",
				"hasRole"
			]
		}
	]
}

```

Then, we pass the `genesis-spec.json` in the following command:

```
burrow configure --genesis-spec=genesis-spec.json > burrow.toml
```

which creates `burrow.toml` that looks like:

```
[GenesisDoc]
  GenesisTime = 2018-07-26T16:33:57Z
  ChainName = "BurrowChain_016017"
  [GenesisDoc.GlobalPermissions]
    Roles = []
    [GenesisDoc.GlobalPermissions.Base]
      Perms = 2302
      SetBit = 16383

  [[GenesisDoc.Accounts]]
    Address = "F71831847564B7008AD30DD56336D9C42787CF63"
    PublicKey = "{\"CurveType\":\"ed25519\",\"PublicKey\":\"1341FFAD471E0AF4CD190293ED0D8D41713439129A0C9C8009A1DD66608B8BF1\"}"
    Amount = 99999999999999
    Name = "Full_0"
    [GenesisDoc.Accounts.Permissions]
      [GenesisDoc.Accounts.Permissions.Base]
        Perms = 16383
        SetBit = 16383

  [[GenesisDoc.Accounts]]
    Address = "8EC3BA4BFE8DDEE76D5C8FF34ED72A83A99DEE98"
    PublicKey = "{\"CurveType\":\"ed25519\",\"PublicKey\":\"E3ACC06A3687E2ECD26818C905B3B73EB54755BF2411C6671AFE522B73FB81C2\"}"
    Amount = 9999999999
    Name = "Participant_0"
    [GenesisDoc.Accounts.Permissions]
      [GenesisDoc.Accounts.Permissions.Base]
        Perms = 2118
        SetBit = 2118

  [[GenesisDoc.Validators]]
    Address = "F71831847564B7008AD30DD56336D9C42787CF63"
    PublicKey = "{\"CurveType\":\"ed25519\",\"PublicKey\":\"1341FFAD471E0AF4CD190293ED0D8D41713439129A0C9C8009A1DD66608B8BF1\"}"
    Amount = 9999999999
    Name = "Full_0"

    [[GenesisDoc.Validators.UnbondTo]]
      Address = "F71831847564B7008AD30DD56336D9C42787CF63"
      PublicKey = "{\"CurveType\":\"ed25519\",\"PublicKey\":\"1341FFAD471E0AF4CD190293ED0D8D41713439129A0C9C8009A1DD66608B8BF1\"}"
      Amount = 9999999999

[Tendermint]
  Seeds = ""
  PersistentPeers = ""
  ListenAddress = "tcp://0.0.0.0:26656"
  ExternalAddress = ""
  Moniker = ""
  TendermintRoot = ".burrow"

[Execution]

[Keys]
  GRPCServiceEnabled = true
  AllowBadFilePermissions = false
  RemoteAddress = ""
  KeysDirectory = ".keys"

[RPC]
  [RPC.TM]
    Enabled = true
    ListenAddress = "tcp://127.0.0.1:26658"
  [RPC.Profiler]
    Enabled = false
    ListenAddress = "tcp://127.0.0.1:6060"
  [RPC.GRPC]
    Enabled = true
    ListenAddress = "127.0.0.1:10997"
  [RPC.Metrics]
    Enabled = false
    ListenAddress = "tcp://127.0.0.1:9102"
    MetricsPath = "/metrics"
    BlockSampleSize = 100

[Logging]
  ExcludeTrace = false
  NonBlocking = false
  [Logging.RootSink]
    [Logging.RootSink.Output]
      OutputType = "stderr"
      Format = "json"

```

## Keys

The previous command (`burrow configure --genesis-spec`) created two keys. Let's look at them:

```
ls .keys/data
```

will show you the existing keys that the `burrow` can use to sign transactions. In this example, signing happens under-the-hood.

## Run Burrow

Now we can run `burrow`:

```
burrow start --validator-index=0 2>burrow.log &
```

See [burrow's README](https://github.com/hyperledger/burrow) for more information on tweaking the logs.

## Deploy Contracts

Now that the burrow node is running, we can deploy contracts.

For this step, we need two things: one or more solidity contracts, and an `epm.yaml`.

Let's take a simple example, found in [this directory](tests/jobs_fixtures/app06-deploy_basic_contract_and_different_solc_types_packed_unpacked/).

The `epm.yaml` should look like:

```
jobs:

- name: deployStorageK
  deploy:
    contract: storage.sol

- name: setStorageBaseBool
  set:
    val: "true"

- name: setStorageBool
  call:
    destination: $deployStorageK
    function: setBool
    data: [$setStorageBaseBool]

- name: queryStorageBool
  query-contract:
    destination: $deployStorageK
    function: getBool

- name: assertStorageBool
  assert:
    key: $queryStorageBool
    relation: eq
    val: $setStorageBaseBool

# tests string bools: #71
- name: setStorageBool2
  call:
    destination: $deployStorageK
    function: setBool2
    data: [true]

- name: queryStorageBool2
  query-contract:
    destination: $deployStorageK
    function: getBool2

- name: assertStorageBool2
  assert:
    key: $queryStorageBool2
    relation: eq
    val: "true"

- name: setStorageBaseInt
  set:
    val: 50000

- name: setStorageInt
  call:
    destination: $deployStorageK
    function: setInt
    data: [$setStorageBaseInt]

- name: queryStorageInt
  query-contract:
    destination: $deployStorageK
    function: getInt

- name: assertStorageInt
  assert:
    key: $queryStorageInt
    relation: eq
    val: $setStorageBaseInt

- name: setStorageBaseUint
  set:
    val: 9999999

- name: setStorageUint
  call:
    destination: $deployStorageK
    function: setUint
    data: [$setStorageBaseUint]

- name: queryStorageUint
  query-contract:
    destination: $deployStorageK
    function: getUint

- name: assertStorageUint
  assert:
    key: $queryStorageUint
    relation: eq
    val: $setStorageBaseUint

- name: setStorageBaseAddress
  set:
    val: "1040E6521541DAB4E7EE57F21226DD17CE9F0FB7"

- name: setStorageAddress
  call:
    destination: $deployStorageK
    function: setAddress
    data: [$setStorageBaseAddress]

- name: queryStorageAddress
  query-contract:
    destination: $deployStorageK
    function: getAddress

- name: assertStorageAddress
  assert:
    key: $queryStorageAddress
    relation: eq
    val: $setStorageBaseAddress

- name: setStorageBaseBytes
  set:
    val: marmatoshi

- name: setStorageBytes
  call:
    destination: $deployStorageK
    function: setBytes
    data: [$setStorageBaseBytes]

- name: queryStorageBytes
  query-contract:
    destination: $deployStorageK
    function: getBytes

- name: assertStorageBytes
  assert:
    key: $queryStorageBytes
    relation: eq
    val: $setStorageBaseBytes

- name: setStorageBaseString
  set:
    val: nakaburrow

- name: setStorageString
  call:
    destination: $deployStorageK
    function: setString
    data: [$setStorageBaseString]

- name: queryStorageString
  query-contract:
    destination: $deployStorageK
    function: getString

- name: assertStorageString
  assert:
    key: $queryStorageString
    relation: eq
    val: $setStorageBaseString
```

while our Solidity contract (`storage.sol`) looks like:

```
pragma solidity >=0.0.0;

contract SimpleStorage {
  bool storedBool;
  bool storedBool2;
  int storedInt;
  uint storedUint;
  address storedAddress;
  bytes32 storedBytes;
  string storedString;

  function setBool(bool x) {
    storedBool = x;
  }

  function getBool() constant returns (bool retBool) {
    return storedBool;
  }

  function setBool2(bool x) {
    storedBool2 = x;
  }

  function getBool2() constant returns (bool retBool) {
    return storedBool2;
  }

  function setInt(int x) {
    storedInt = x;
  }

  function getInt() constant returns (int retInt) {
    return storedInt;
  }

  function setUint(uint x) {
    storedUint = x;
  }

  function getUint() constant returns (uint retUint) {
    return storedUint;
  }

  function setAddress(address x) {
    storedAddress = x;
  }

  function getAddress() constant returns (address retAddress) {
    return storedAddress;
  }

  function setBytes(bytes32 x) {
    storedBytes = x;
  }

  function getBytes() constant returns (bytes32 retBytes) {
    return storedBytes;
  }

  function setString(string x) {
    storedString = x;
  }

  function getString() constant returns (string retString) {
    return storedString;
  }
}
```

Both files (`epm.yaml` & `storage.sol`) should be in the same directory with no other yaml or sol files.

From inside that directory, we are ready to deploy.

```
bos --address=F71831847564B7008AD30DD56336D9C42787CF63
```

where you should replace the `--address` field with the `ValidatorAddress` at the top of your `burrow.toml`.

That's it! You've succesfully deployed (and tested) a Soldity contract to a Burrow node.

Note - that to redeploy the burrow chain later, you will need the same genesis-spec.json and burrow.toml files - so keep hold of them!

## Working with Javascript

Currently the javascript libraries are being rebuilt. The master branch of this repository works against the master branch of burrow.

Please use the versions within the package.json of this repo on master branch for fully compatible and tested versions.
