# Bosmarmot

Bosmarmot is a monorepo containing condensed and updated versions of the monax tooling. This repo intends to provide the basic tooling required to interact with a [Burrow](https://github.com/hyperledger/burrow) chain.

It also contains the interpreter for the Monax packages specification language 'epm'. This README will cover setting up a Burrow chain with the Monax tooling from start to finish.

## Install

We're going to need four (4) binaries:


```
burrow
bos
monax-keys
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

For `bos` and `monax-keys`:

```
go get github.com/monax/bosmarmot
cd $GOPATH/src/github.com/monax/bosmarmot
make build
```

and move these onto you `$PATH` as well.

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
	Accounts": [
		{
			"Amount": 99999999999999,
			"AmountBonded": 9999999999,
			"Name": "Full_0",
			"Permissions": [
				"all"
			]
		},
		{
			"Amount": 9999999999,
			"Name": "Participant_0",
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

Because the next command will be making keys, let's start keys server:

```
monax-keys server 2>keys.log &
```

Then, we pass the `genesis-spec.json` in the following command:

```
burrow configure --genesis-spec=genesis-spec.json > burrow.toml
```

which creates `burrow.toml` that looks like:

```
[GenesisDoc]
  GenesisTime = 2018-05-24T16:12:34Z
  ChainName = "BurrowChain_2BE507"
  [GenesisDoc.GlobalPermissions]
    Roles = []
    [GenesisDoc.GlobalPermissions.Base]
      Perms = 2302
      SetBit = 16383

  [[GenesisDoc.Accounts]]
    Address = "84276E34B8F3C4166F61878473AEC831A5CF5444"
    PublicKey = "{\"type\":\"ed25519\",\"data\":\"8C245576A2A0299DF0B760C55EF366D281B914FADBE03571FC1901FCAADA067F\"}"
    Amount = 99999999999999
    Name = "Full_0"
    [GenesisDoc.Accounts.Permissions]
      [GenesisDoc.Accounts.Permissions.Base]
        Perms = 16383
        SetBit = 16383

  [[GenesisDoc.Accounts]]
    Address = "8F573B59FE9D0886CB525069E7E16C1997E8A126"
    PublicKey = "{\"type\":\"ed25519\",\"data\":\"15989611C2670C248F81070F7E94E824B527E6CA80C4F4EE05414C82F9DA18E3\"}"
    Amount = 9999999999
    Name = "Participant_0"
    [GenesisDoc.Accounts.Permissions]
      [GenesisDoc.Accounts.Permissions.Base]
        Perms = 2118
        SetBit = 2118

  [[GenesisDoc.Validators]]
    Address = "84276E34B8F3C4166F61878473AEC831A5CF5444"
    PublicKey = "{\"type\":\"ed25519\",\"data\":\"8C245576A2A0299DF0B760C55EF366D281B914FADBE03571FC1901FCAADA067F\"}"
    Amount = 9999999999
    Name = "Full_0"

    [[GenesisDoc.Validators.UnbondTo]]
      Address = "84276E34B8F3C4166F61878473AEC831A5CF5444"
      PublicKey = "{\"type\":\"ed25519\",\"data\":\"8C245576A2A0299DF0B760C55EF366D281B914FADBE03571FC1901FCAADA067F\"}"
      Amount = 9999999999

[Tendermint]
  Seeds = ""
  PersistentPeers = ""
  ListenAddress = "tcp://0.0.0.0:46656"
  Moniker = ""
  TendermintRoot = ".burrow"

[Keys]
  URL = "http://localhost:4767"

[RPC]
  [RPC.V0]
    Disabled = false
    [RPC.V0.Server]
      [RPC.V0.Server.bind]
        address = "localhost"
        port = 1337
      [RPC.V0.Server.TLS]
        tls = false
        cert_path = ""
        key_path = ""
      [RPC.V0.Server.CORS]
        enable = false
        allow_credentials = false
        max_age = 0
      [RPC.V0.Server.HTTP]
        json_rpc_endpoint = "/rpc"
      [RPC.V0.Server.web_socket]
        websocket_endpoint = "/socketrpc"
        max_websocket_sessions = 50
        read_buffer_size = 4096
        write_buffer_size = 4096
  [RPC.TM]
    Disabled = false
    ListenAddress = "tcp://localhost:46657"
  [RPC.Profiler]
    Disabled = true
    ListenAddress = "tcp://localhost:6060"

[Logging]
  ExcludeTrace = false
  NonBlocking = false
  [Logging.RootSink]
    [Logging.RootSink.Output]
      OutputType = "stderr"
      Format = "json"

```

## Keys

Recall that we ran `monax-keys server` to start the keys server. The previous command (`burrow configure --genesis-spec`) created two keys. Let's look at them:

```
ls $HOME/.monax/keys/data
```

will show you the existing keys that the `monax-keys server` can use to sign transactions. In this example, signing happens under-the-hood.

## Run Burrow

Now we can run `burrow`:

```
burrow start --validator-index=0 2>burrow.log &
```

See [burrow's README](https://github.com/hyperledger/burrow) for more information on tweaking the logs.

## Deploy Contracts

Now that the burrow node is running, we can deploy contracts.

For this step, we need two things: one or more solidity contracts, and an `epm.yaml`.

Let's take a simple example, found in [this directory](monax/tests/jobs_fixtures/app06-deploy_basic_contract_and_different_solc_types_packed_unpacked/).

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
      data:
        - $setStorageBaseBool

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
      data:
        - true

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
      data:
        - $setStorageBaseInt

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
      data:
        - $setStorageBaseUint

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
      data:
        - $setStorageBaseAddress

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
      data:
        - $setStorageBaseBytes

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
      data:
        - $setStorageBaseString

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

Both files (`epm.yaml` & `storage.sol`) should be in the same directory with nothing else in it.

From inside that directory, we are ready to deploy.

```
bos pkgs do --keys="http://localhost:4767" --chain-url="tcp://localhost:46657" --address=0A40DC874BC932B78AC390EAD1C1BF33469597AB
```

where the field in `--address` is the `ValidatorAddress` at the top of your `burrow.toml`.

That's it! You've succesfully deployed (and tested) a Soldity contract to a Burrow node.
