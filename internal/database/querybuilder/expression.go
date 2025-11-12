package querybuilder

import (
	"fmt"
	"strings"
)

// Expr represents a SQL expression that can be converted to SQL string
type Expr interface {
	ToSQL() (string, []any)
}

// ColumnExpr represents a column reference in SQL
type ColumnExpr struct {
	Name string
}

// ToSQL converts column expression to SQL string
func (c *ColumnExpr) ToSQL() (string, []any) {
	return c.Name, nil
}

// LiteralExpr represents a literal value in SQL (will be converted to placeholder)
type LiteralExpr struct {
	Value any
}

// ToSQL converts literal value to SQL placeholder
func (l *LiteralExpr) ToSQL() (string, []any) {
	return "?", []any{l.Value}
}

// UnaryExpr represents a unary operation (IS NULL, IS NOT NULL, NOT)
type UnaryExpr struct {
	Operator Op
	Operand  Expr
}

// ToSQL converts unary expression to SQL string
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

// TernaryExpr represents a ternary operation (BETWEEN)
type TernaryExpr struct {
	Expr Expr
	Low  Expr
	High Expr
}

// ToSQL converts ternary expression (BETWEEN) to SQL string
func (b *TernaryExpr) ToSQL() (string, []any) {
	colSQL, colArgs := b.Expr.ToSQL()
	lowSQL, lowArgs := b.Low.ToSQL()
	highSQL, highArgs := b.High.ToSQL()

	sql := fmt.Sprintf("%s BETWEEN %s AND %s", colSQL, lowSQL, highSQL)
	args := append(colArgs, lowArgs...)
	args = append(args, highArgs...)

	return sql, args
}

// BinaryExpr represents a binary operation (=, !=, >, <, AND, OR, etc.)
type BinaryExpr struct {
	Operator Op
	Left     Expr
	Right    Expr
}

// ToSQL converts binary expression to SQL string using in-order traversal
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
