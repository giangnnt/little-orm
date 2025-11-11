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
		if err := (*v).ValidateAndTransform(&expr.Left); err != nil {
			return err
		}
		if err := (*v).ValidateAndTransform(&expr.Right); err != nil {
			return err
		}
	case *LiteralExpr:
		return nil
	default:
	}
	return nil
}
