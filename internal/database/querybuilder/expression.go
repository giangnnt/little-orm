package querybuilder

import (
	"fmt"
	"strings"
)

type Expr interface {
	ToSQL() (string, []any)
}

type ColumnExpr struct {
	Name string
}

func (c *ColumnExpr) ToSQL() (string, []any) {
	return c.Name, nil
}

type LiteralExpr struct {
	Value any
}

func (l *LiteralExpr) ToSQL() (string, []any) {
	return "?", []any{l.Value}
}

type BinaryExpr struct {
	Operator Op
	Left     Expr
	Right    Expr
}

func (b *BinaryExpr) ToSQL() (string, []any) {
	leftSQL, leftArgs := b.Left.ToSQL()
	rightSQL, rightArgs := b.Right.ToSQL()

	var sql strings.Builder
	var args []any

	switch b.Operator {
	case OpAnd, OpOr, OpIn, OpNIn:
		sql.WriteString(fmt.Sprintf("(%s %s %s)", leftSQL, b.Operator, rightSQL))
		args = append(leftArgs, rightArgs...)
	case OpEq, OpNEq, OpGt, OpLt, OpGte, OpLte, OpLike:
		sql.WriteString(fmt.Sprintf("%s %s %s", leftSQL, b.Operator, rightSQL))
		args = append(leftArgs, rightArgs...)
	default:
		return "", nil
	}

	return sql.String(), args
}
