'use strict'

const assert = require('assert')
const test = require('../../lib/test')

const Test = test.Test()

describe('#44', function () {
  before(Test.before())
  after(Test.after())

  this.timeout(60 * 1000)

  it('#44', Test.it(function (burrow) {
    const source = `
      contract SimpleStorage {
          address storedData;

          function set(address x) {
              storedData = x;
          }

          function get() constant returns (address retVal) {
              return storedData;
          }
      }
    `
    const {abi, bytecode} = test.compile(source, 'SimpleStorage')
    return burrow.contracts.deploy(abi, bytecode).then((contract) =>
      contract.set('88977A37D05A4FE86D09E88C88A49C2FCF7D6D8F')
        .then(() => contract.get())
    ).then((value) => {
      assert.equal(value, '88977A37D05A4FE86D09E88C88A49C2FCF7D6D8F')
    })
  }))
})
