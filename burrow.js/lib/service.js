'use strict'
const path = require('path')

const PROTO_PATH = path.join(__dirname, '../protobuf/')
const PROTO_FILE = 'exec.proto'

const protoLoader = require('@grpc/proto-loader')
const grpc = require('grpc')

const options = {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
  includeDirs: [
    PROTO_PATH
  ]
}

const packageDefinition = protoLoader.loadSync(PROTO_FILE, options)

function removeNested (object) {
  if (!object) return

  var newObject = {}

  for (var key in object) {
    if (key === 'nested') {
      return removeNested(object.nested)
    }

    if (typeof object[key] === 'object' && object[key].constructor === Object) {
      newObject[key] = removeNested(object[key])
    } else {
      newObject[key] = object[key]
    }
  }
  return newObject
}

function wrapGRPC (name) {
  return function (params, callback) {
    // Fetch requestType and ResponseType
    var pName = this.packageName
    var sName = this.serviceName

    var reqStream = this.pbJSON[pName][sName].methods[name].requestStream
    var resStream = this.pbJSON[pName][sName].methods[name].responseStream

    if (reqStream) {
      throw new Error("Can't call a requestStream method")
    }

    if (resStream) {
      if (!callback) throw new Error('Callback not provided')

      var call = this.client[name](params)
      call.on('data', (data) => {
        callback(null, data)
      })
      call.on('error', (error) => {
        callback(error)
      })
      // Return a function that will close the stream when called
      return () => {
        console.log('WARNING: stream closing is not implemented')
      }
    } else {
      // Make call through client
      var P = new Promise((resolve, reject) => {
        this.client[name](params, function (err, result) {
          if (err) return reject(err)
          resolve(result)
        })
      })
      if (callback) {
        return P.then((result) => { callback(null, result) }).catch((err) => { callback(err) })
      } else {
        return P
      }
    }
  }
}

function Service (filePath, packageName, serviceName, URL) {
  this.packageName = packageName
  this.serviceName = serviceName

  this.service = grpc.loadPackageDefinition(packageDefinition)
  this.pbJSON = removeNested(packageDefinition)

  this.client = new this.service[packageName][serviceName](URL, grpc.credentials.createInsecure())

  for (var method in this.pbJSON[packageName][serviceName].methods) {
    this[method] = wrapGRPC(method).bind(this)
  }
}

module.exports = function (file, packageName, serviceName, URL) {
  return new Service(file, packageName, serviceName, URL)
}
