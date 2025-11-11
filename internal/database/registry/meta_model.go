package registry

type ColumnMeta struct {
	DBTag string
	Name  string
	Type  string
	Tag   string
}

type TableMeta struct {
	TableName string
	Columns   map[string]ColumnMeta
}

func (t TableMeta) HasColumn(columnName string) bool {
	_, ok := t.Columns[columnName]
	return ok
}
