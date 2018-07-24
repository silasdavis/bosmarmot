# @monax/legacy-db.js (Alpha)

This is a JavaScript API for communicating with a [Hyperledger Burrow](https://github.com/hyperledger/burrow) server.

## New Name

This library used to be named `eris-db.js`.  It is now `@monax/legacy-db.js` as part of the company-wide renaming to Monax and also to distinguish it from the upcoming new client API.  Although it is a legacy API it will continue to be supported.

To use new versions of the library in existing code, change the line in your `package.json` which looks like this:

```
"eris-db": "0.15.12",
```

to make it look like this:

```
"@monax/legacy-db": "0.16.0",
```

and run `npm install`.

## Installation

### Prerequisites

* [Git](https://git-scm.com/)
* [Monax](https://monax.io/) version 0.16
* [Node.js](https://nodejs.org/) version 6 or higher

You can check the installed version of Node.js with the command:

```shell
$ node --version
```

If your distribution of Linux has a version older than 6 then you can
update it using [NodeSource's distribution](https://github.com/nodesource/distributions).

### To Install

```shell
$ npm install @monax/legacy-db
```

## Usage

If you created a Burrow server using the [Monax CLI](https://github.com/monax/cli) tool, you can find out its IP address using the following command:

```
$ monax chains ip <name of Burrow server>
```

The main class is `Burrow`. A standard `Burrow` instance is created like this:

```JavaScript
var burrowFactory = require('@monax/legacy-db');

var burrow = burrowFactory.createInstance("http://<IP address>:1337/rpc");
```

The parameters for `createInstance` is the server URL as a string. The client-type is chosen based on the URL scheme. As of now, the supported schemes are: `http` and `ws` (websockets). No additional configuration is needed.

## API Reference

There are bindings for all the RPC methods. All functions are on the form `function(param1, param2, ... , callback)`, where the callback is a function on the form `function(error,data)` (it is documented under the name `methodCallback`). The `data` object is the same as you would get by calling the corresponding RPC method directly.

This is the over-all structure of the library. The `unsafe` flag means a private key is either sent or received, so should be used with care (dev only).

NOTE: There will be links to the proper jsdoc and integration with Monax.io. For now, the components point to the actual code files and methods points to the web-API method in question.

### Burrow

| Component Name | Accessor |
| :------------- | :------- |
| Accounts | [Burrow.accounts()](https://github.com/monax/legacy-db.js/blob/master/lib/accounts.js) |
| Blockchain | [Burrow.blockchain()](https://github.com/monax/legacy-db.js/blob/master/lib/blockchain.js) |
| Consensus | [Burrow.consensus()](https://github.com/monax/legacy-db.js/blob/master/lib/consensus.js) |
| Events | [Burrow.events()](https://github.com/monax/legacy-db.js/blob/master/lib/events.js) |
| NameReg | [Burrow.namereg()](https://github.com/monax/legacy-db.js/blob/master/lib/namereg.js) |
| Network | [Burrow.network()](https://github.com/monax/legacy-db.js/blob/master/lib/network.js) |
| Transactions | [Burrow.txs()](https://github.com/monax/legacy-db.js/blob/master/lib/transactions.js) |

### Components

#### Accounts

The accounts object has methods for getting account and account-storage data.

| Method | RPC method | Notes |
| :----- | :--------- | :---- |
| Accounts.getAccounts | [Burrow.getAccounts](https://monax.io/docs/documentation/db/latest/specifications/api/#getaccounts) | |
| Accounts.getAccount | [Burrow.getAccount](https://monax.io/docs/documentation/db/latest/specifications/api/#getaccount) | |
| Accounts.getStorage | [Burrow.getStorage](https://monax.io/docs/documentation/db/latest/specifications/api/#getstorage) | |
| Accounts.getStorageAt | [Burrow.getStorageAt](https://monax.io/docs/documentation/db/latest/specifications/api/#getstorageat) | |
| Accounts.genPrivAccount | [Burrow.genPrivAccount](https://monax.io/docs/documentation/db/latest/specifications/api/#genprivaccount) | unsafe |

#### BlockChain

The accounts object has methods for getting blockchain-related data, such as a list of blocks, or individual blocks, or the hash of the genesis block.

| Method | RPC method | Notes |
| :----- | :--------- | :---- |
| BlockChain.getInfo |  [Burrow.getBlockchainInfo](https://monax.io/docs/documentation/db/latest/specifications/api/#getblockchaininfo) | |
| BlockChain.getChainId | [Burrow.getChainId](https://monax.io/docs/documentation/db/latest/specifications/api/#getchainid) | |
| BlockChain.getGenesisHash | [Burrow.getGenesisHash](https://monax.io/docs/documentation/db/latest/specifications/api/#getgenesishash) | |
| BlockChain.getLatestBlockHeight | [Burrow.getLatestBlockHeight](https://monax.io/docs/documentation/db/latest/specifications/api/#getlatestblockheight) | |
| BlockChain.getLatestBlock | [Burrow.getLatestBlock](https://monax.io/docs/documentation/db/latest/specifications/api/#getlatestblock) | |
| BlockChain.getBlocks | [Burrow.getBlocks](https://monax.io/docs/documentation/db/latest/specifications/api/#getblocks) | |
| BlockChain.getBlock | [Burrow.getBlock](https://monax.io/docs/documentation/db/latest/specifications/api/#getblock) | |

#### Consensus

The consensus object has methods for getting consensus-related data.

| Method | RPC method | Notes |
| :----- | :--------- | :---- |
| Consensus.getState |   [Burrow.getConsensusState](https://monax.io/docs/documentation/db/latest/specifications/api/#getconsensusstate) | |
| Consensus.getValidators | [Burrow.getValidators](https://monax.io/docs/documentation/db/latest/specifications/api/#getvalidators) | |

#### Events

The tendermint client will generate and fire off events when important things happen, like when a new block has been committed, or someone is transacting to an account. It is possible to subscribe to these events. These are the methods for subscribing, un-subscribing and polling.

| Method | RPC method | Notes |
| :----- | :--------- | :---- |
| Events.subscribe | [Burrow.eventSubscribe](https://monax.io/docs/documentation/db/latest/specifications/api/#eventsubscribe) | |
| Events.unsubscribe | [Burrow.eventUnsubscribe](https://monax.io/docs/documentation/db/latest/specifications/api/#eventunubscribe) | |
| Events.poll | [Burrow.eventPoll](https://monax.io/docs/documentation/db/latest/specifications/api/#eventpoll) | |

##### Helpers

The helper functions makes it easier to manage subscriptions. Normally you'd be using these functions rather then managing the subscriptions yourself.

Helper functions always contain two callback functions - a `createCallback(error, data)` and an `eventCallback(error, data)`.

The `createCallback` data is an [EventSub]() object, that can be used to do things like getting the event ID, the subscriber ID, and to stop the subscription.

The `eventCallback` data is the event object. This object is different depending on the event type. In the case of `NewBlock` it will be a block, the consensus events is a transaction object, etc. More info can be found in the [api doc]().

| Method | Arguments |
| :----- | :-------- |
| Events.subAccountInput | `account address <string>` |
| Events.subAccountOutput | `account address <string>` |
| Events.subAccountReceive | `account address <string>` |
| Events.subLogEvent | `account address <string>` |
| Events.subSolidityEvent | `account address <string>` |
| Events.subNewBlocks | `-` |
| Events.subForks | `-` |
| Events.subBonds | `-` |
| Events.subUnbonds | `-` |
| Events.subRebonds | `-` |
| Events.subDupeouts | `-` |

`subSolidityEvent` and `subLogEvent` are two different names for the same type of subscription (log events).

#### NameReg

The NameReg object has methods for accessing the name registry.

| Method | RPC method | Notes |
| :----- | :--------- | :---- |
| NameReg.getEntry | [Burrow.getNameRegEntry](https://monax.io/docs/documentation/db/latest/specifications/api/#get-namereg-entry) | |
| NameReg.getEntries | [Burrow.getNameRegEntries](https://monax.io/docs/documentation/db/latest/specifications/api/#get-namereg-entries) | |

#### Network

The accounts object has methods for getting network-related data, such as a list of all peers. It could also have been named "node".

Client Version may be a bit misplaced

| Method | RPC method | Notes |
| :----- | :--------- | :---- |
| Network.getInfo | [Burrow.getNetworkInfo](https://monax.io/docs/documentation/db/latest/specifications/api/#getnetworkinfo) |  |
| Network.getClientVersion | [Burrow.getClientVersion](https://monax.io/docs/documentation/db/latest/specifications/api/#getclientversion) | |
| Network.getMoniker | [Burrow.getMoniker](https://monax.io/docs/documentation/db/latest/specifications/api/#getmoniker) | |
| Network.isListening | [Burrow.isListening](https://monax.io/docs/documentation/db/latest/specifications/api/#islistening) | |
| Network.getListeners | [Burrow.getListeners](https://monax.io/docs/documentation/db/latest/specifications/api/#getlisteners) | |
| Network.getPeers | [Burrow.getPeers](https://monax.io/docs/documentation/db/latest/specifications/api/#getpeers) | |
| Network.getPeer | [Burrow.getPeer](https://monax.io/docs/documentation/db/latest/specifications/api/#getpeer) | |

#### Transactions

A transaction is the equivalence of a database `write` operation. They can be done in two ways. There's the "dev" way, which is to call `transact` and pass along the target address (if any), data, gas, and a private key used for signing. It is very similar to the old Ethereum way of transacting, except Tendermint does not keep accounts in the client, so a private key needs to be sent along. This means the server **should either run on the same machine as the tendermint client, or in the same, private network**.

Transacting via `broadcastTx` will be the standard way of doing things if you want the key to remain on the users machine. This requires a browser plugin for doing the actual signing, which we will add later. For now, you should stick to the `transact` method.

To get a private key for testing/developing, you can run `tendermint gen_account` if you have it installed. You can also run `tools/pa_generator.js` if you have a local node running. It will take the url as command line argument at some point...

##### Calls

Calls provide read-only access to the smart contracts. It is used mostly to get data out of a contract-accounts storage by using the contracts accessor methods, but can be used to call any method that does not change any data in any account. A trivial example would be a contract function that takes two numbers as input, adds them, and then simply returns the sum.

There are two types of calls. `Call` takes a data string and an account address and calls the code in that account (if any) using the provided data as input. This is the standard method for read-only operations.

`CallCode` works the same except you don't provide an account address but the actual compiled code instead. It's a dev tool for accessing the VM directly. "Code-execution as a service".

| Method | RPC method | Notes |
| :----- | :--------- | :---- |
| Transactions.broadcastTx | [Burrow.broadcastTx](https://monax.io/docs/documentation/db/latest/specifications/api/#broadcasttx) | see below |
| Transactions.getUnconfirmedTxs | [Burrow.getUnconfirmedTxs](https://monax.io/docs/documentation/db/latest/specifications/api/#getunconfirmedtxs) | |
| Transactions.call | [Burrow.call](https://monax.io/docs/documentation/db/latest/specifications/api/#call) | |
| Transactions.callCode | [Burrow.callCode](https://monax.io/docs/documentation/db/latest/specifications/api/#callcode) | |
| Transactions.transact | [Burrow.transact](https://monax.io/docs/documentation/db/latest/specifications/api/#transact) | unsafe |
| Transactions.transactAndHold | [Burrow.transactAndHold](https://monax.io/docs/documentation/db/latest/specifications/api/#transact-and-hold) | unsafe |
| Transactions.transactNameReg | [Burrow.transactNameReg](https://monax.io/docs/documentation/db/latest/specifications/api/#transactnamereg) | unsafe |

`broadcastTx` is useless until we add a client-side signing solution.

## Documentation

Generate documentation using the command `npm run doc`.

## Testing

To test the library against pre-recorded vectors:

```
npm test
```

To test the library against Burrow while recording vectors:

```
TEST=record npm test
```

To test Burrow against pre-recorded vectors without exercising the library:

```
TEST=server npm test
```

## Debugging

Debugging information will display on `stderr` if the library is run with `NODE_DEBUG=monax` in the environment.

## Copyright

Copyright 2015 Monax

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
