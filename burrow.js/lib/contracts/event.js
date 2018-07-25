var utils = require('../utils/utils')
var coder = require('web3/lib/solidity/coder')
var sha3 = require('../utils/sha3')

/**
 * This prototype should be used to create event filters
 */

var types = function (abi, indexed) {
  return abi.inputs.filter(function (i) {
    return i.indexed === indexed
  }).map(function (i) {
    return i.type
  })
}
/**
 * Should be used to decode indexed params and options
 *
 * @method decode
 * @param {Object} data
 * @return {Object} result object with decoded indexed && not indexed params
 */
var decode = function (abi, data) {
  data.Data = data.Data.toString('hex').toUpperCase() || ''
  data.Topics = data.Topics.map(x => x.toString('hex').toUpperCase()) || []

  var argTopics = abi.anonymous ? data.Topics : data.Topics.slice(1)
  var indexedData = argTopics.join('')
  var indexedParams = coder.decodeParams(types(abi, true), indexedData)

  // var notIndexedData = data.Data.slice(2)
  var notIndexedParams = coder.decodeParams(types(abi, false), data.Data)

  // var result = formatters.outputLogFormatter(data);
  var result = {}
  result.event = utils.transformToFullName(abi)
  result.address = data.Address

  result.args = abi.inputs.reduce(function (acc, current) {
    acc[current.name] = current.indexed ? indexedParams.shift() : notIndexedParams.shift()
    return acc
  }, {})

  return result
}

var SolidityEvent = function (abi) {
  var name = utils.transformToFullName(abi)
  var displayName = utils.extractDisplayName(name)
  var typeName = utils.extractTypeName(name)
  var signature = sha3(name)

  var call = function (address, callback) {
    address = address || this.address
    if (!callback) { throw new Error('Can not subscribe to an event without a callback') };

    return this.burrow.pipe.eventSub(address, signature, (error, event) => {
      if (error) return callback(error)

      try {
        var decoded = decode(abi, event.Log)
      } catch (error) {
        return callback(error)
      }

      const converted = Object.assign(
        {},
        decoded,
        {args: utils.web3ToBurrow(decoded.args)}
      )

      return callback(null, converted)
    })
  }

  return {displayName, typeName, call}
}

module.exports = SolidityEvent
