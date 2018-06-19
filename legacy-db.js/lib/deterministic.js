// A proxy that makes Burrow's communication appear deterministic in order to
// ease testing.

'use strict'

const generic = require('@nodeguy/generic')
const _ = generic._
const I = require('iteray')
const {is} = require('@nodeguy/type')
const R = require('ramda')

module.exports = (transport) => {
  let subIds = new Map()

  const newRequest = generic.function()

  newRequest.method(_, R.identity)

  newRequest.method(
    [{
      params: {
        subId: is(String)
      }
    }],
    (request) =>
      R.assocPath(
        ['params', 'subId'],
        subIds.get(request.params.subId),
        request
      )
  )

  const newResponse = generic.function()

  newResponse.method(_, R.identity)

  // events
  //
  // Replace each non-deterministic event 'height' with '1'.
  newResponse.method(
    [{
      result: {
        events: is(Array)
      }
    }],
    (response) =>
      R.assocPath(['result', 'events'], response.result.events.map((event) => {
        event.Execution.Log.Height = 1
        return event
      }
      ), response)
  )

  // sub_id
  newResponse.method(
    [{
      result: {
        subId: is(String)
      }
    }],
    (response) => {
      const newId = String(subIds.size)
      subIds.set(newId, response.result.subId)
      return R.assocPath(['result', 'subId'], newId, response)
    }
  )

  // tx_hash
  //
  // Override the nondeterministic 'tx_hash' with the deterministic JSON-RPC id.
  newResponse.method([{
    id: is(String),
    result: {
      TxHash: is(String)
    }
  }], (response) =>
    R.assocPath(['result', 'TxHash'], response.id, response)
  )

  return (input) =>
    I.map(newResponse, transport(I.map(newRequest, input)))
}
