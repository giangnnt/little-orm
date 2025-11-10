package querybuilder

type SQLBuilderType string

type QueryString string

const (
	SelectType SQLBuilderType = "select"
	InsertType SQLBuilderType = "insert"
)

type SortOrder string

const (
	Ascending  SortOrder = "ASC"
	Descending SortOrder = "DESC"
)

type LogicOp string

const (
	AND LogicOp = "AND"
	OR  LogicOp = "OR"
)
