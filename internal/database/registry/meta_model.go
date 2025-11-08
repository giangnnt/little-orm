package registry

type ColumnMeta struct {
	Name string
	Type string
	Tag  string
}

type TableMeta struct {
	TableName string
	Columns   []ColumnMeta
}
