package querybuilder

import (
	"little-orm/internal/database/registry"
)

type InsertBuilder struct {
	table     string
	columns   []string
	values    []any
	tableMeta registry.TableMeta
}

func NewInsertBuilder(model any) *InsertBuilder {
	// Get table registry and table meta
	reg := registry.GetDBRegistry()
	tableMeta := reg.GetTableMeta(model)

	return &InsertBuilder{
		tableMeta: tableMeta,
		table:     tableMeta.TableName,
		columns:   make([]string, 0),
		values:    make([]any, 0),
	}
}

func (b *InsertBuilder) Build() (string, []any) {
	// TODO: Implement build logic
	return "", nil
}
