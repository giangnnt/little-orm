package querybuilder

type Expr struct {
	Field    string
	Operator string
	Value    any
	Logic    LogicOp
}
