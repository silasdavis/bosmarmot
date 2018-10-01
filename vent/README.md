# Vent Component

Vent reads a sqlsol json events specification file, parses its contents and maps column types to corresponding sql types to synchronize database structures.

Then listens to burrow gRPC events based on given filters, parses, decodes data and builds rows to be upserted in corresponding event tables.

Rows are upserted in blocks, where each block for a given event filter is one commit.

Block id and event filtering data is stored in Log tables in order to resume pending blocks.

Given a sqlsol specification 

```json
[
	{
		"TableName" : "TEST_TABLE",
		"Filter" : "Log1Text = 'EVENT_TEST'",
		"Event"  : {
			"anonymous": false,
			"inputs": [{
				"indexed": true,
				"name": "name",
				"type": "bytes32"
			}, {
				"indexed": false,
				"name": "key",
				"type": "uint256"
			}, {
				"indexed": false,
				"name": "blocknum",
				"type": "uint256"
			}, {
				"indexed": false,
				"name": "somestr",
				"type": "string"
			}, {
				"indexed": false,
				"name": "this",
				"type": "address"
			}, {
				"indexed": false,
				"name": "instance",
				"type": "uint256"
			}],
			"name": "UpdateTable",
			"type": "event"
		},
		"Columns"  : {
			"key"		: {"name" : "Index",    "primary" : true},
			"blocknum"  : {"name" : "Block",    "primary" : false},
			"somestr"	: {"name" : "String",   "primary" : false},
			"instance" 	: {"name" : "Instance", "primary" : false}
		}
	}
]
```

Vent builds dictionary, log and event tables for the defined tables & columns and maps input types to proper sql types.

Database structures are created or altered on the fly based on specifications (just adding new columns is supported).


## Adapters:

Adapters are database implementations, Vent can store data in different rdbms.

In `sqldb/adapters` there's a list of supported adapters (there is also a README.md file in that folder that helps to understand how to implement a new one).

## Setup PostgreSQL Database with Docker:

```bash
# Create postgres container (only once):
docker run --name postgres-local -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=vent -p 5432:5432 -d postgres:10.4-alpine

# Start postgres container:
docker start postgres-local

# Stop postgres container:
docker stop postgres-local

# Delete postgres container:
docker container rm postgres-local
```

## Run Unit Tests:

```bash
# From the main repo folder:
make test_integration_vent
```

## Run Vent Command:

```bash
# Install vent command:
go install ./vent

# Print command help:
vent --help

# Run vent command with postgres adapter and spec file path:
vent --db-adapter="postgres" --db-url="postgres://user:pass@localhost:5432/vent?sslmode=disable" --db-schema="vent" --grpc-addr="localhost:10997" --http-addr="0.0.0.0:8080" --log-level="debug" --spec-file="<sqlsol specification file path>"

# Run vent command with sqlite adapter and spec directory path:
vent --db-adapter="sqlite" --db-url="./vent.sqlite" --grpc-addr="localhost:10997" --http-addr="0.0.0.0:8080" --log-level="debug" --spec-dir="<sqlsol specification directory path>"
```

Configuration Flags:

+ `db-adapter`: Database adapter, 'postgres' or 'sqlite' are fully supported
+ `db-url`: PostgreSQL database URL or SQLite db file path
+ `db-schema`: PostgreSQL database schema or empty for SQLite
+ `http-addr`: Address to bind the HTTP server
+ `grpc-addr`: Address to listen to gRPC Hyperledger Burrow server
+ `log-level`: Logging level (error, warn, info, debug)
+ `spec-file`: SQLSol specification json file (full path)
+ `spec-dir`: Path of a folder to look for SQLSol json specification files


NOTES:

One of `spec-file` or `spec-dir` must be provided.
If `spec-dir` is provided, Vent will search for all `.json` spec files in given directory.
It can be checked that Vent is connected and ready sending a request to `http://<http-addr>/health` which will return a `200` OK response in case everything's fine.
