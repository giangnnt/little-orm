package querybuilder

import (
	"fmt"
	"little-orm/internal/database/registry"
)

type ExprValidator struct {
	tableMeta registry.TableMeta
}

func (v *ExprValidator) ValidateAndTransform(expr *Expr) error {
	switch expr := (*expr).(type) {
	case *ColumnExpr:
		colMeta, ok := v.tableMeta.Columns[expr.Name]
		if !ok {
			return fmt.Errorf("column '%s' not found in table '%s'", expr.Name, v.tableMeta.TableName)
		}
		expr.Name = colMeta.DBTag
	case *BinaryExpr:
		if expr.Left != nil {
			if err := (*v).ValidateAndTransform(&expr.Left); err != nil {
				return err
			}
		}
		if expr.Right != nil {
			if err := (*v).ValidateAndTransform(&expr.Right); err != nil {
				return err
			}
		}
	case *UnaryExpr:
		if expr.Operand != nil {
			if err := (*v).ValidateAndTransform(&expr.Operand); err != nil {
				return err
			}
		}
	case *TernaryExpr:
		if expr.Expr != nil {
			if err := (*v).ValidateAndTransform(&expr.Expr); err != nil {
				return err
			}
		}
		if expr.Low != nil {
			if err := (*v).ValidateAndTransform(&expr.Low); err != nil {
				return err
			}
		}
		if expr.High != nil {
			if err := (*v).ValidateAndTransform(&expr.High); err != nil {
				return err
			}
		}
	case *LiteralExpr:
		return nil
	default:
	}
	return nil
}
