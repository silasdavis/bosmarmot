package test

import (
	"testing"
)

// GoodJSONConfFile sets a good json file to be used in parser tests
func GoodJSONConfFile(t *testing.T) string {
	t.Helper()

	goodJSONConfFile := `[
		{
			"TableName" : "UserAccounts",
			"Filter" : "LOG0 = 'UserAccounts'",
			"Event"  : {
				"anonymous": false,
				"inputs": [{
					"indexed": false,
					"name": "userName",
					"type": "string"
				}, {
					"indexed": false,
					"name": "userAddress",
					"type": "address"
				}, {
					"indexed": false,
					"name": "UnimportantInfo",
					"type": "uint"
				}],
				"name": "UpdateUserAccount",
				"type": "event"
			},
			"Columns"  : {
				"userAddress" : {"name" : "address", "primary" : true},
				"userName": {"name" : "username", "primary" : false}
			}
		},
		{
			"TableName" : "EventTest",
			"Filter" : "LOG0 = 'EventTest'",
			"Event"  : {
				"anonymous": false,
				"inputs": [{
					"indexed": false,
					"name": "name",
					"type": "string"
				}, {
					"indexed": false,
					"name": "description",
					"type": "string"
				}, {
					"indexed": false,
					"name": "UnimportantInfo",
					"type": "uint"
				}],
				"name": "TEST_EVENTS",
				"type": "event"
			},
			"Columns"  : {
				"name" : {"name" : "testname", "primary" : true},
				"description": {"name" : "testdescription", "primary" : false}
			}
		}
	]`

	return goodJSONConfFile
}

// MissingFieldsJSONConfFile sets a json file with missing fields to be used in parser tests
func MissingFieldsJSONConfFile(t *testing.T) string {
	t.Helper()

	missingFieldsJSONConfFile := `[
		{
			"TableName" : "UserAccounts",
			"Event"  : {
				"anonymous": false,
				"inputs": [{
					"indexed": false,
					"name": "userName",
					"type": "string"
				}, {
					"indexed": false,
					"name": "userAddress",
					"type": "address"
				}, {
					"indexed": false,
					"name": "UnimportantInfo",
					"type": "uint"
				}],
				"type": "event"
			},
			"Columns"  : {
				"userAddress" : {"name" : "address", "primary" : true},
				"userName": {"name" : "username", "primary" : false}
			}
		}
	]`

	return missingFieldsJSONConfFile
}

// UnknownTypeJSONConfFile sets a json file with unknown column types to be used in parser tests
func UnknownTypeJSONConfFile(t *testing.T) string {
	t.Helper()

	unknownTypeJSONConfFile := `[
		{
			"TableName" : "UserAccounts",
			"Filter" : "LOG0 = 'UserAccounts'",
			"Event"  : {
				"anonymous": false,
				"inputs": [{
					"indexed": false,
					"name": "userName",
					"type": "typeunknown"
				}, {
					"indexed": false,
					"name": "userAddress",
					"type": "address"
				}, {
					"indexed": false,
					"name": "UnimportantInfo",
					"type": "uint"
				}],
				"name": "UpdateUserAccount",
				"type": "event"
			},
			"Columns"  : {
				"userAddress" : {"name" : "address", "primary" : true},
				"userName": {"name" : "username", "primary" : false}
			}
		},
		{
			"TableName" : "EventTest",
			"Filter" : "LOG0 = 'EventTest'",
			"Event"  : {
				"anonymous": false,
				"inputs": [{
					"indexed": false,
					"name": "name",
					"type": "typeunknown"
				}, {
					"indexed": false,
					"name": "description",
					"type": "string"
				}, {
					"indexed": false,
					"name": "UnimportantInfo",
					"type": "uint"
				}],
				"name": "TEST_EVENTS",
				"type": "event"
			},
			"Columns"  : {
				"name" : {"name" : "testname", "primary" : true},
				"description": {"name" : "testdescription", "primary" : false}
			}
		}
	]`

	return unknownTypeJSONConfFile
}

// BadJSONConfFile sets a malformed json file to be used in parser tests
func BadJSONConfFile(t *testing.T) string {
	t.Helper()

	badJSONConfFile := `[
		{
			"TableName" : "UserAccounts",
			"Event"  : {
				"anonymous": false,
				"inputs": [{
					"indexed": false,
					"name": "userName",
					"type": "string"
				}, {
					"indexed": false,
					"name": "userAddress",
					"type": "address"
				}, {
					"indexed": false,
					"name": "UnimportantInfo",
					"type": "uint"
				}],
				"name": "UpdateUserAccount",
			},
			"Columns"  : {
				"userAddress" : {"name" : "address", "primary" : true},
				"userName": {"name" : "username", "primary" : false}
	]`

	return badJSONConfFile
}
