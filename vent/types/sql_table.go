package types

// SQLTable contains the structure of a SQL table,
type SQLTable struct {
	Name    string
	Filter  string
	Columns map[string]SQLTableColumn
}

// SQLTableColumn contains the definition of a SQL table column,
// the Order is given to be able to sort the columns to be created
type SQLTableColumn struct {
	Name    string
	Type    SQLColumnType
	Length  int
	Primary bool
	Order   int
}

// UpsertQuery contains generic query to upsert row data
type UpsertQuery struct {
	Query   string
	Length  int
	Columns map[string]UpsertColumn
}

// UpsertColumn contains info about a specific column to be upserted
type UpsertColumn struct {
	IsNumeric   bool
	InsPosition int
	UpdPosition int
}

// defined fixed sql column names in log and event tables
const (
	SQLColumnNameHeight      = "_height"
	SQLColumnNameIndex       = "_index"
	SQLColumnNameTxHash      = "_txhash"
	SQLColumnNameEventName   = "_eventname"
	SQLColumnNameEventType   = "_eventtype"
	SQLColumnNameEventFilter = "_eventfilter"
	SQLColumnNameTableName   = "_tablename"
	SQLColumnNameTimeStamp   = "_timestamp"
	SQLColumnNameId          = "_id"
	SQLColumnNameRowCount    = "_rowcount"
)

// SQL log table
const SQLLogTableName = "_vent_log"
