/**
 * @file dev_pipe.js
 * @fileOverview Base class for the dev-pipe.
 * @author Andreas Olofsson
 * @module pipe/dev_pipe
 */
'use strict'

var Pipe = require('./pipe')
var nUtil = require('util')
var config = require('../utils/config')

var ZERO_ADDRESS = '0000000000000000000000000000000000000000'

/**
 * Constructor for the DevPipe class.
 *
 * @type {Pipe}
 */
module.exports = DevPipe

/**
 * DevPipe transacts using the unsafe private-key transactions.
 *
 * @param {*} burrow - the burrow object.
 * @param {string} accounts - the private key to use when sending transactions. NOTE: This means a private key
 * will be passed over the net, so it should only be used when developing, or if it's 100% certain that the
 * Burrow server and this script runs on the same machine, or communication is secure. The recommended way
 * will be to call a signing function on the client side, like in a browser plugin.
 *
 * @constructor
 */
function DevPipe (burrow, account, options) {
  Pipe.call(this, burrow)

  // Choose signing mode
  this._options = Object.assign({signbyaddress: false, readonly: false}, options)

  if (!account) {
    this._options.readonly = true
  } else {
    this._account = _formatAccount(account, this._options.signbyaddress)
  }
}

nUtil.inherits(DevPipe, Pipe)

/**
 * Used to send a transaction.
 * @param {module:solidity/function~TxPayload} txPayload - The payload object.
 * @param callback - The error-first callback. The 'data' param is a contract address in the case of a
 * create transactions, otherwise it's the return value.
 */
DevPipe.prototype.transact = function (txPayload, callback) {
  var to = txPayload.to
  var from

  if (txPayload.from) {
    try {
      from = _formatAccount(txPayload.from)
    } catch (err) {
      return callback(new Error('No account matches the provided address/ID: ' + txPayload.from))
    }
  } else {
    if (this._options.readonly) {
      return callback(new Error('Pipe is readonly and no account address provided'))
    }

    from = this._account
  }
  this._burrow.txs().transactAndHold(from, to, txPayload.data, config.DEFAULT_GAS, config.DEFAULT_FEE, function (error, data) {
    if (error) {
      console.log(error)
      callback(error)
    } else {
      if (to) {
        callback(null, data.Return.toUpperCase())
      } else {
        callback(null, data.CallData.Callee)
      }
    }
  })
}

/**
 * Used to do a call.
 * @param {module:solidity/function~TxPayload} txPayload - The payload object.
 * @param callback - The error-first callback.
 */
DevPipe.prototype.call = function (txPayload, callback) {
  var address = txPayload.to
  var from = txPayload.from
  if (!from) {
    from = ZERO_ADDRESS
  }
  var data = txPayload.data
  if (data.length > 1 && data.slice(0, 2) === '0x') {
    data = data.slice(2)
  }
  this._burrow.txs().call(from, address, data, function (error, data) {
    if (error) {
      callback(error)
    } else {
      callback(null, data.Return.toUpperCase())
    }
  })
}

function _formatAccount (account, signbyaddress) {
  var ad = {
    address: undefined,
    privateKey: undefined
  }

  if (typeof (account) === 'string') {
    if (signbyaddress) {
      ad.address = account
    } else {
      ad.privateKey = account
    }
  } else if (typeof (account) === 'object') {
    if (!account.address && !account.privKey) {
      throw new Error('Account data is not on the proper format: ' + JSON.stringify(account))
    }

    if (signbyaddress && account.address != null) {
      ad.address = account.address
    } else if (account.privKey != null) {
      ad.privateKey = account.privKey
    } else {
      throw new Error('Account object does not provide sufficient data: ' + JSON.stringify(account))
    }
  } else {
    throw new Error('Account data is not in the proper format: ' + JSON.stringify(account))
  }

  return ad
}
