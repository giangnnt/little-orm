package querybuilder

// Basic helper
func C(name string) Expr             { return &ColumnExpr{Name: name} }
func L(val any) Expr                 { return &LiteralExpr{Value: val} }
func U(op Op, operand Expr) Expr     { return &UnaryExpr{Operator: op, Operand: operand} }
func B(op Op, left, right Expr) Expr { return &BinaryExpr{Operator: op, Left: left, Right: right} }
func T(expr, low, high Expr) Expr    { return &TernaryExpr{Expr: expr, Low: low, High: high} }

// Logical helper
func And(exprs ...Expr) Expr {
	if len(exprs) == 0 {
		return nil
	}
	e := exprs[0]
	for _, ex := range exprs[1:] {
		e = B(OpAnd, e, ex)
	}
	return e
}

func Or(exprs ...Expr) Expr {
	if len(exprs) == 0 {
		return nil
	}
	e := exprs[0]
	for _, ex := range exprs[1:] {
		e = B(OpOr, e, ex)
	}
	return e
}
