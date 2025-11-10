package querybuilder

type InsertBuilder struct {
	table   string
	columns []string
	values  []any
}
