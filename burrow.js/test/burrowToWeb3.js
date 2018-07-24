'use strict'

const assert = require('assert')
const utils = require('../lib/utils/utils')

const monax = '000000000000000000000000000000000000000000000000000000000000002A'

const web3 =
  '0x000000000000000000000000000000000000000000000000000000000000002A'

it('converts a Monax hex string to a Web3 hex value', function () {
  assert.equal(utils.burrowToWeb3(monax), web3)
})

it('converts a Web3 hex value to a Monax hex string', function () {
  assert.equal(utils.web3ToBurrow(web3), monax)
})
