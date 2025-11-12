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

type UnaryExpr struct {
	Operator Op
	Operand  Expr
}

func (u *UnaryExpr) ToSQL() (string, []any) {
	operandSQL, args := u.Operand.ToSQL()
	switch u.Operator {
	case "IS NULL", "IS NOT NULL":
		return fmt.Sprintf("%s %s", operandSQL, u.Operator), args
	case "NOT":
		return fmt.Sprintf("NOT (%s)", operandSQL), args
	default:
		panic("unsupported unary operator: " + u.Operator)
	}
}

type TernaryExpr struct {
	Expr Expr
	Low  Expr
	High Expr
}

func (b *TernaryExpr) ToSQL() (string, []any) {
	colSQL, colArgs := b.Expr.ToSQL()
	lowSQL, lowArgs := b.Low.ToSQL()
	highSQL, highArgs := b.High.ToSQL()

	sql := fmt.Sprintf("%s BETWEEN %s AND %s", colSQL, lowSQL, highSQL)
	args := append(colArgs, lowArgs...)
	args = append(args, highArgs...)

	return sql, args
}

type BinaryExpr struct {
	Operator Op
	Left     Expr
	Right    Expr
}

// In-order traversal
func (b *BinaryExpr) ToSQL() (string, []any) {
	// Check for nil operands
	if b.Left == nil || b.Right == nil {
		panic("Binary Expression don't have enough Left or Right")
	}

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
