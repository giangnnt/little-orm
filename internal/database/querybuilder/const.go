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

type Op string

const (
	// Comparison operators
	OpEq        Op = "="
	OpNEq       Op = "!="
	OpGt        Op = ">"
	OpLt        Op = "<"
	OpGte       Op = ">="
	OpLte       Op = "<="
	OpLike      Op = "LIKE"
	OpIn        Op = "IN"
	OpNIn     Op = "NOT IN"
	OpIsNull    Op = "IS NULL"
	OpIsNotNull Op = "IS NOT NULL"
	OpBetween   Op = "BETWEEN"

	// Logical operators
	OpAnd Op = "AND"
	OpOr  Op = "OR"
	OpNot Op = "NOT"
)
