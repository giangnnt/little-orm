package querybuilder

import (
	"fmt"
	"little-orm/internal/database/registry"
	"strings"
)

// SelectBuilder builds SELECT SQL queries
type SelectBuilder struct {
	table         string
	fields        []Expr
	exprs         Expr
	groupBy       []ColumnExpr
	orderBy       []string
	sortOrder     []SortOrder
	limit         int
	offset        int
	args          []any
	tableMeta     registry.TableMeta
	exprValidator *ExprValidator
}

// NewSelectBuilder creates a new SELECT query builder for the given model
func NewSelectBuilder(model any) *SelectBuilder {
	// Get table registry and table meta
	reg := registry.GetDBRegistry()
	tableMeta := reg.GetTableMeta(model)

	// Init all fields
	fields := make([]Expr, 0, len(tableMeta.Columns))
	for _, col := range tableMeta.Columns {
		fields = append(fields, C(col.DBTag))
	}

	return &SelectBuilder{
		tableMeta:     tableMeta,
		table:         tableMeta.TableName,
		fields:        fields,
		exprValidator: &ExprValidator{tableMeta: tableMeta},
	}
}

// Select specifies which fields to select (if not called, selects all fields)
func (b *SelectBuilder) Select(fields ...string) *SelectBuilder {
	dbTags := []Expr{}
	for _, field := range fields {
		fieldMeta, ok := b.tableMeta.Columns[field]
		if !ok {
			panic(fmt.Sprintf("Field %s is not registered", field))
		}
		dbTags = append(dbTags, C(fieldMeta.DBTag))
	}
	// Replace init fields
	b.fields = dbTags
	return b
}

// Where adds WHERE clause to the query
func (b *SelectBuilder) Where(e Expr) *SelectBuilder {
	if err := b.exprValidator.ValidateAndTransform(&e); err != nil {
		panic(err.Error())
	}
	b.exprs = e
	return b
}

// OrderBy adds ORDER BY clause to the query
func (b *SelectBuilder) OrderBy(order string, sortOrder SortOrder) *SelectBuilder {
	b.orderBy = append(b.orderBy, order)
	b.sortOrder = append(b.sortOrder, sortOrder)
	return b
}

// Limit sets the LIMIT clause
func (b *SelectBuilder) Limit(n int) *SelectBuilder {
	b.limit = n
	return b
}

// Offset sets the OFFSET clause
func (b *SelectBuilder) Offset(m int) *SelectBuilder {
	b.offset = m
	return b
}

// Build constructs the final SQL query and returns it with arguments
func (b *SelectBuilder) Build() (string, []any) {
	query := b.buildSelectClause()
	query += b.buildWhereClause()
	query += b.buildOrderByClause()
	query += b.buildLimitOffsetClause()
	return query, b.args
}

// buildWhereClause constructs the WHERE clause
func (b *SelectBuilder) buildWhereClause() string {
	if b.exprs == nil {
		return ""
	}
	whereClause, args := b.exprs.ToSQL()
	b.args = args
	return " WHERE " + whereClause
}

// buildSelectClause constructs the SELECT clause
func (b *SelectBuilder) buildSelectClause() string {
	fields := "*"
	// Giả sử b.fields có kiểu []*ColumnExpr
	expr, ok := (b.fields).([]*ColumnExpr)
	if ok {
		names := make([]string, 0, len(expr))
		for _, v := range expr {
			names = append(names, v.Name)
		}
		fieldsStr := ""
		if len(expr) > 0 {
			fieldsStr = strings.Join(names, ", ")
		}
		return fmt.Sprintf("SELECT %s FROM %s", fieldsStr, b.table)
	}

	// Nếu không phải []*ColumnExpr, có thể xử lý khác hoặc trả về lỗi
	return ""

}

// buildOrderByClause constructs the ORDER BY clause
func (b *SelectBuilder) buildOrderByClause() string {
	if len(b.orderBy) == 0 {
		return ""
	}

	orders := make([]string, len(b.orderBy))
	for i := range b.orderBy {
		ord := b.orderBy[i]
		if i < len(b.sortOrder) {
			ord += " " + string(b.sortOrder[i])
		}
		orders[i] = ord
	}

	return " ORDER BY " + strings.Join(orders, ", ")
}

// buildLimitOffsetClause constructs the LIMIT and OFFSET clauses
func (b *SelectBuilder) buildLimitOffsetClause() string {
	result := ""
	if b.limit > 0 {
		result += fmt.Sprintf(" LIMIT %d", b.limit)
	}
	if b.offset > 0 {
		result += fmt.Sprintf(" OFFSET %d", b.offset)
	}
	return result
}
