package querybuilder

import (
	"fmt"
	"little-orm/internal/database/registry"
	"strings"
)

type SelectBuilder struct {
	table      string
	fields     []string
	conditions []string
	orderBy    []string
	sortOrder  []SortOrder
	limit      int
	offset     int
	args       []any
	tableMeta  registry.TableMeta
}

func NewSelectBuilder(model any) *SelectBuilder {
	// Get table registry and table meta
	reg := registry.GetDBRegistry()
	tableMeta := reg.GetTableMeta(model)

	// Init all fields
	fields := make([]string, 0, len(tableMeta.Columns))
	for _, col := range tableMeta.Columns {
		fields = append(fields, col.DBTag)
	}

	return &SelectBuilder{
		tableMeta: tableMeta,
		fields:    fields,
	}
}

func (b *SelectBuilder) Select(fields ...string) *SelectBuilder {
	b.fields = fields
	return b
}

func (b *SelectBuilder) Where(cond string, args ...any) *SelectBuilder {
	b.conditions = append(b.conditions, cond)
	b.args = append(b.args, args...)
	return b
}

func (b *SelectBuilder) OrderBy(order string, sortOrder SortOrder) *SelectBuilder {
	b.orderBy = append(b.orderBy, order)
	b.sortOrder = append(b.sortOrder, sortOrder)
	return b
}

func (b *SelectBuilder) Limit(n int) *SelectBuilder {
	b.limit = n
	return b
}

func (b *SelectBuilder) Offset(m int) *SelectBuilder {
	b.offset = m
	return b
}

func (b *SelectBuilder) Build() (string, []any) {
	// SELECT clause
	fields := ""
	if len(b.fields) > 0 {
		fields = strings.Join(b.fields, ", ")
	} else {
		fields = "*"
	}
	query := fmt.Sprintf("SELECT %s FROM %s", fields, b.table)

	// WHERE clause
	if len(b.conditions) > 0 {
		query += " WHERE " + strings.Join(b.conditions, " AND ")
	}

	// ORDER BY clause
	if len(b.orderBy) > 0 {
		orders := make([]string, len(b.orderBy))
		for i := range b.orderBy {
			ord := b.orderBy[i]
			if i < len(b.sortOrder) {
				ord += " " + string(b.sortOrder[i])
			}
			orders[i] = ord
		}
		query += " ORDER BY " + strings.Join(orders, ", ")
	}

	// LIMIT / OFFSET
	if b.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", b.limit)
	}
	if b.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", b.offset)
	}

	return query, b.args
}
