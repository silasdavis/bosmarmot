package types

import (
	"github.com/go-ozzo/ozzo-validation"
)

// EventSpec contains all event specifications
type EventSpec []EventDefinition

// EventDefinition struct (table name where to persist filtered events and it structure)
type EventDefinition struct {
	TableName string                 `json:"TableName"`
	Filter    string                 `json:"Filter"`
	Columns   map[string]EventColumn `json:"Columns"`
}

// Validate checks the structure of an EventDefinition
func (evDef EventDefinition) Validate() error {
	return validation.ValidateStruct(&evDef,
		validation.Field(&evDef.TableName, validation.Required, validation.Length(1, 60)),
		validation.Field(&evDef.Filter, validation.Required),
		validation.Field(&evDef.Columns, validation.Required, validation.Length(1, 0)),
	)
}

// EventColumn struct (table column definition)
type EventColumn struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Primary     bool   `json:"primary"`
	HexToString bool   `json:"hexToString"`
}

// Validate checks the structure of an EventColumn
func (evColumn EventColumn) Validate() error {
	return validation.ValidateStruct(&evColumn,
		validation.Field(&evColumn.Name, validation.Required, validation.Length(1, 60)),
	)
}
