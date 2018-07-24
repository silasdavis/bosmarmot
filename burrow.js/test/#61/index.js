'use strict'

const assert = require('assert')
const test = require('../../lib/test')

const Test = test.Test()

describe('#61', function () {
  before(Test.before())
  after(Test.after())

  this.timeout(10 * 1000)

  it('#61', Test.it(function (burrow) {
    const source = `
      contract SimpleStorage {
          address storedData;

          function SimpleStorage(address x) {
              storedData = x;
          }

          function get() constant returns (address retVal) {
              return storedData;
          }
      }
    `
    const {abi, bytecode} = test.compile(source, 'SimpleStorage')
    return burrow.contracts.deploy(abi, bytecode, '88977A37D05A4FE86D09E88C88A49C2FCF7D6D8F')
      .then((contract) => contract.get())
      .then((value) => {
        assert.equal(value, '88977A37D05A4FE86D09E88C88A49C2FCF7D6D8F')
      })
  }))
})
