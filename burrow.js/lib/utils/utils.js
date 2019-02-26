/**
 * @file utils.js
 * @author Marek Kotewicz <marek@ethdev.com>
 * @author Andreas Olofsson
 * @date 2015
 * @module utils/utils
 */

const coder = require('ethereumjs-abi')
const sha3 = require('./sha3')

/**
 * Should be used to create full function/event name from json abi
 *
 * @method transformToFullName
 * @param {Object} json - json-abi
 * @return {String} full fnction/event name
 */
var transformToFullName = function (json) {
  if (json.name.indexOf('(') !== -1) {
    return json.name
  }

  var typeName = json.inputs.map(function (i) {
    return i.type
  }).join()
  return json.name + '(' + typeName + ')'
}

/**
 * Should be called to get display name of contract function
 *
 * @method extractDisplayName
 * @param {String} name of function/event
 * @returns {String} display name for function/event eg. multiply(uint256) -> multiply
 */
var extractDisplayName = function (name) {
  var length = name.indexOf('(')
  return length !== -1 ? name.substr(0, length) : name
}

/**
 *
 * @param {String} name - the name.
 * @returns {String} overloaded part of function/event name
 */
var extractTypeName = function (name) {
  /// TODO: make it invulnerable
  var length = name.indexOf('(')
  return length !== -1 ? name.substr(length + 1, name.length - 1 - (length + 1)).replace(' ', '') : ''
}

/**
 * Auto converts any given value into it's hex representation.
 *
 * @method toHex
 * @param {boolean|String|Number|BigNumber|Object} val - the value.
 * @return {String}
 */

/**
 * Returns true if object is function, otherwise false
 *
 * @method isFunction
 * @param {Object} object - object to test
 * @return {Boolean}
 */
var isFunction = function (object) {
  return typeof object === 'function'
}

var encode = function (abi, functionName, args) {
  var functions = abi.filter(function (json) {
    return (json.type === 'function' && json.name === functionName)
  })

  if (functions.length === 0) {
    throw new Error('Function name: ' + functionName + ' not found in abi')
  } else if (functions.length > 1) {
    throw new Error('Function name: ' + functionName + ' is overloaded, Overloading is not supported')
  } else {
    var name = transformToFullName(functions[0])
    var functionSig = sha3(name).slice(0, 8)
    var types = functions[0].inputs.map(function (arg) {
      return arg.type
    })
    return functionSig + coder.rawEncode(types, args)
  }
}

module.exports = {
  transformToFullName: transformToFullName,
  extractDisplayName: extractDisplayName,
  extractTypeName: extractTypeName,
  isFunction: isFunction,
  encode: encode
}
