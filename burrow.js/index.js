/**
 * @file index.js
 * @fileOverview Index file for the Burrow javascript API. This file contains a factory method
 * for creating a new <tt>Burrow</tt> instance.
 * @author Andreas Olofsson
 * @module index
 */
'use strict'

var Burrow = require('./lib/Burrow')
var utils = require('./lib/utils/utils')

module.exports = {
  createInstance: Burrow.createInstance,
  Burrow: Burrow
}

/**
 * Utils has methods for working with strings.
 *
 * @type {{}}
 */
exports.utils = {}
exports.utils.hexToAscii = utils.hexToAscii
exports.utils.asciiToHex = utils.asciiToHex
exports.utils.padLeft = utils.padLeft
exports.utils.padRight = utils.padRight
exports.utils.htoa = utils.hexToAscii
exports.utils.atoh = utils.asciiToHex
