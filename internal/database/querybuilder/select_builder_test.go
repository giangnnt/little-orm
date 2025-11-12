package querybuilder

import (
	"fmt"
	"little-orm/internal/model"
	"strings"
	"testing"
)

func TestNewSelectBuilder(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})

	if builder == nil {
		t.Fatal("Expected builder to not be nil")
	}

	if builder.table != "users" {
		t.Errorf("Expected table to be 'users', got '%s'", builder.table)
	}

	if len(builder.fields) == 0 {
		t.Error("Expected fields to be initialized with columns")
	}
}

func TestSelectBuilder_Build_Simple(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Build()

	// Check query structure (field order may vary due to map iteration)
	if !strings.Contains(query, "SELECT") {
		t.Error("Expected query to contain SELECT")
	}
	if !strings.Contains(query, "FROM users") {
		t.Error("Expected query to contain FROM users")
	}
	// Check all fields are present
	if !strings.Contains(query, "id") || !strings.Contains(query, "email") ||
		!strings.Contains(query, "name") || !strings.Contains(query, "password") {
		t.Errorf("Expected query to contain all fields (id, email, name, password), got: %s", query)
	}

	if len(args) != 0 {
		t.Errorf("Expected no args, got %d", len(args))
	}
}

func TestSelectBuilder_Where(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpEq,
		Left:     &ColumnExpr{Name: "Email"},
		Right:    &LiteralExpr{Value: "test@example.com"},
	}).Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "WHERE email = ?") {
		t.Errorf("Unexpected query: %s", query)
	}

	if len(args) != 1 || args[0] != "test@example.com" {
		t.Errorf("Expected args ['test@example.com'], got %v", args)
	}
}

func TestSelectBuilder_Where_Multiple(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpAnd,
		Left: &BinaryExpr{
			Operator: OpEq,
			Left:     &ColumnExpr{Name: "Email"},
			Right:    &LiteralExpr{Value: "test@example.com"},
		},
		Right: &BinaryExpr{
			Operator: OpEq,
			Left:     &ColumnExpr{Name: "Name"},
			Right:    &LiteralExpr{Value: "John"},
		},
	}).Build()

	// Check query structure (field order may vary due to map iteration)
	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "WHERE (email = ? AND name = ?)") {
		t.Errorf("Unexpected query: %s", query)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}

	if args[0] != "test@example.com" || args[1] != "John" {
		t.Errorf("Expected args ['test@example.com', 'John'], got %v", args)
	}
}

func TestSelectBuilder_OrderBy_Ascending(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, _ := builder.OrderBy("name", Ascending).Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "ORDER BY name ASC") {
		t.Errorf("Unexpected query: %s", query)
	}
}

func TestSelectBuilder_OrderBy_Descending(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, _ := builder.OrderBy("id", Descending).Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "ORDER BY id DESC") {
		t.Errorf("Unexpected query: %s", query)
	}
}

func TestSelectBuilder_OrderBy_Multiple(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, _ := builder.
		OrderBy("name", Ascending).
		OrderBy("id", Descending).
		Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "ORDER BY name ASC, id DESC") {
		t.Errorf("Unexpected query: %s", query)
	}
}

func TestSelectBuilder_Limit(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, _ := builder.Limit(10).Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "LIMIT 10") {
		t.Errorf("Unexpected query: %s", query)
	}
}

func TestSelectBuilder_Offset(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, _ := builder.Offset(5).Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "OFFSET 5") {
		t.Errorf("Unexpected query: %s", query)
	}
}

func TestSelectBuilder_LimitAndOffset(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, _ := builder.Limit(10).Offset(20).Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "LIMIT 10") || !strings.Contains(query, "OFFSET 20") {
		t.Errorf("Unexpected query: %s", query)
	}
}

func TestSelectBuilder_ComplexQuery(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.
		OrderBy("name", Ascending).
		Where(&BinaryExpr{
			Operator: OpAnd,
			Left: &BinaryExpr{
				Operator: OpAnd,
				Left: &BinaryExpr{
					Operator: OpLike,
					Left:     &ColumnExpr{Name: "Email"},
					Right:    &LiteralExpr{Value: "%@example.com"},
				},
				Right: &BinaryExpr{
					Operator: OpNEq,
					Left:     &ColumnExpr{Name: "Name"},
					Right:    &LiteralExpr{Value: "Admin"},
				},
			},
			Right: &UnaryExpr{
				Operator: OpIsNNull,
				Operand:  &ColumnExpr{Name: "ID"},
			},
		}).
		OrderBy("id", Descending).
		Limit(25).
		Offset(50).
		Build()

	// Check all parts are present (field order may vary due to map iteration)
	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "WHERE ((email LIKE ? AND name != ?) AND id IS NOT NULL)") ||
		!strings.Contains(query, "ORDER BY name ASC, id DESC") ||
		!strings.Contains(query, "LIMIT 25") || !strings.Contains(query, "OFFSET 50") {
		t.Errorf("Unexpected query: %s", query)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}

	if args[0] != "%@example.com" || args[1] != "Admin" {
		t.Errorf("Unexpected args: %v", args)
	}
}

func TestSelectBuilder_Select_InvalidField(t *testing.T) {
	setupTestRegistry()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid field, but didn't panic")
		}
	}()

	builder := NewSelectBuilder(model.User{})
	builder.Select("NonExistentField")
}

func TestSelectBuilder_Chaining(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})

	// Test that methods return *SelectBuilder for chaining
	result := builder.Where(&BinaryExpr{
		Operator: OpGt,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    &LiteralExpr{Value: 1},
	})
	if result != builder {
		t.Error("Where should return the same builder instance")
	}

	result = builder.OrderBy("name", Ascending)
	if result != builder {
		t.Error("OrderBy should return the same builder instance")
	}

	result = builder.Limit(10)
	if result != builder {
		t.Error("Limit should return the same builder instance")
	}

	result = builder.Offset(5)
	if result != builder {
		t.Error("Offset should return the same builder instance")
	}
}

func TestSelectBuilder_EmptyConditions(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, _ := builder.Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") {
		t.Errorf("Unexpected query: %s", query)
	}
	// Should not contain WHERE, ORDER BY, LIMIT, or OFFSET
	if strings.Contains(query, "WHERE") || strings.Contains(query, "ORDER BY") ||
		strings.Contains(query, "LIMIT") || strings.Contains(query, "OFFSET") {
		t.Errorf("Expected simple query without conditions, got: %s", query)
	}
}

// Test all comparison operators
func TestSelectBuilder_Where_OpEq(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpEq,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    &LiteralExpr{Value: 1},
	}).Build()

	if !strings.Contains(query, "WHERE id = ?") {
		t.Errorf("Expected WHERE clause with id = ?, got: %s", query)
	}

	if len(args) != 1 || args[0] != 1 {
		t.Errorf("Expected args [1], got %v", args)
	}
}

func TestSelectBuilder_Where_OpNEq(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpNEq,
		Left:     &ColumnExpr{Name: "Name"},
		Right:    &LiteralExpr{Value: "Admin"},
	}).Build()

	if !strings.Contains(query, "WHERE name != ?") {
		t.Errorf("Expected WHERE clause with name != ?, got: %s", query)
	}

	if len(args) != 1 || args[0] != "Admin" {
		t.Errorf("Expected args ['Admin'], got %v", args)
	}
}

func TestSelectBuilder_Where_OpGt(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpGt,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    &LiteralExpr{Value: 10},
	}).Build()

	if !strings.Contains(query, "WHERE id > ?") {
		t.Errorf("Expected WHERE clause with id > ?, got: %s", query)
	}

	if len(args) != 1 || args[0] != 10 {
		t.Errorf("Expected args [10], got %v", args)
	}
}

func TestSelectBuilder_Where_OpLt(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpLt,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    &LiteralExpr{Value: 100},
	}).Build()

	if !strings.Contains(query, "WHERE id < ?") {
		t.Errorf("Expected WHERE clause with id < ?, got: %s", query)
	}

	if len(args) != 1 || args[0] != 100 {
		t.Errorf("Expected args [100], got %v", args)
	}
}

func TestSelectBuilder_Where_OpGte(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpGte,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    &LiteralExpr{Value: 5},
	}).Build()

	if !strings.Contains(query, "WHERE id >= ?") {
		t.Errorf("Expected WHERE clause with id >= ?, got: %s", query)
	}

	if len(args) != 1 || args[0] != 5 {
		t.Errorf("Expected args [5], got %v", args)
	}
}

func TestSelectBuilder_Where_OpLte(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpLte,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    &LiteralExpr{Value: 50},
	}).Build()

	if !strings.Contains(query, "WHERE id <= ?") {
		t.Errorf("Expected WHERE clause with id <= ?, got: %s", query)
	}

	if len(args) != 1 || args[0] != 50 {
		t.Errorf("Expected args [50], got %v", args)
	}
}

func TestSelectBuilder_Where_OpLike(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpLike,
		Left:     &ColumnExpr{Name: "Email"},
		Right:    &LiteralExpr{Value: "%@gmail.com"},
	}).Build()

	if !strings.Contains(query, "WHERE email LIKE ?") {
		t.Errorf("Expected WHERE clause with email LIKE ?, got: %s", query)
	}

	if len(args) != 1 || args[0] != "%@gmail.com" {
		t.Errorf("Expected args ['%%@gmail.com'], got %v", args)
	}
}

func TestSelectBuilder_Where_OpIn(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpIn,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    &LiteralExpr{Value: []int{1, 2, 3}},
	}).Build()

	if !strings.Contains(query, "WHERE (id IN ?)") {
		t.Errorf("Expected WHERE clause with (id IN ?), got: %s", query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestSelectBuilder_Where_OpNIn(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpNIn,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    &LiteralExpr{Value: []int{4, 5, 6}},
	}).Build()

	if !strings.Contains(query, "WHERE (id NOT IN ?)") {
		t.Errorf("Expected WHERE clause with (id NOT IN ?), got: %s", query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

// Test logical operators
func TestSelectBuilder_Where_OpOr(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpOr,
		Left: &BinaryExpr{
			Operator: OpEq,
			Left:     &ColumnExpr{Name: "Name"},
			Right:    &LiteralExpr{Value: "John"},
		},
		Right: &BinaryExpr{
			Operator: OpEq,
			Left:     &ColumnExpr{Name: "Name"},
			Right:    &LiteralExpr{Value: "Jane"},
		},
	}).Build()

	if !strings.Contains(query, "WHERE (name = ? OR name = ?)") {
		t.Errorf("Expected WHERE clause with (name = ? OR name = ?), got: %s", query)
	}

	if len(args) != 2 || args[0] != "John" || args[1] != "Jane" {
		t.Errorf("Expected args ['John', 'Jane'], got %v", args)
	}
}

// Test complex nested expressions
func TestSelectBuilder_Where_ComplexNested(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	// (email LIKE '%@gmail.com' OR email LIKE '%@yahoo.com') AND id > 10
	query, args := builder.Where(&BinaryExpr{
		Operator: OpAnd,
		Left: &BinaryExpr{
			Operator: OpOr,
			Left: &BinaryExpr{
				Operator: OpLike,
				Left:     &ColumnExpr{Name: "Email"},
				Right:    &LiteralExpr{Value: "%@gmail.com"},
			},
			Right: &BinaryExpr{
				Operator: OpLike,
				Left:     &ColumnExpr{Name: "Email"},
				Right:    &LiteralExpr{Value: "%@yahoo.com"},
			},
		},
		Right: &BinaryExpr{
			Operator: OpGt,
			Left:     &ColumnExpr{Name: "ID"},
			Right:    &LiteralExpr{Value: 10},
		},
	}).Build()

	if !strings.Contains(query, "WHERE ((email LIKE ? OR email LIKE ?) AND id > ?)") {
		t.Errorf("Expected WHERE clause with nested conditions, got: %s", query)
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}

	if args[0] != "%@gmail.com" || args[1] != "%@yahoo.com" || args[2] != 10 {
		t.Errorf("Expected args ['%%@gmail.com', '%%@yahoo.com', 10], got %v", args)
	}
}

// Test Select specific fields
func TestSelectBuilder_Select_SpecificFields(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	builder.fields = []string{} // Clear default fields
	query, _ := builder.Select("ID", "Email").Build()

	expectedQuery := "SELECT id, email FROM users"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}
}

// Test invalid column in expression
func TestSelectBuilder_Where_InvalidColumn(t *testing.T) {
	setupTestRegistry()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid column, but didn't panic")
		}
	}()

	builder := NewSelectBuilder(model.User{})
	builder.Where(&BinaryExpr{
		Operator: OpEq,
		Left:     &ColumnExpr{Name: "NonExistentColumn"},
		Right:    &LiteralExpr{Value: "value"},
	})
}

// Test NULL operators (IS NULL, IS NOT NULL) using UnaryExpr
func TestSelectBuilder_Where_OpIsNull(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&UnaryExpr{
		Operator: OpIsNull,
		Operand:  &ColumnExpr{Name: "Email"},
	}).Build()

	if !strings.Contains(query, "WHERE email IS NULL") {
		t.Errorf("Expected WHERE clause with email IS NULL, got: %s", query)
	}

	if len(args) != 0 {
		t.Errorf("Expected 0 args for IS NULL, got %d", len(args))
	}
}

func TestSelectBuilder_Where_OpIsNotNull(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&UnaryExpr{
		Operator: OpIsNNull,
		Operand:  &ColumnExpr{Name: "Name"},
	}).Build()

	if !strings.Contains(query, "WHERE name IS NOT NULL") {
		t.Errorf("Expected WHERE clause with name IS NOT NULL, got: %s", query)
	}

	if len(args) != 0 {
		t.Errorf("Expected 0 args for IS NOT NULL, got %d", len(args))
	}
}

// Test panic when right is nil for operators that require it
func TestSelectBuilder_Where_OpEq_NilRight_ShouldPanic(t *testing.T) {
	setupTestRegistry()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for OpEq with nil Right, but didn't panic")
		}
	}()

	builder := NewSelectBuilder(model.User{})

	// Build should panic when Right is nil for binary expression
	builder.Where(&BinaryExpr{
		Operator: OpEq,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    nil, // This is invalid for OpEq
	}).Build()
}

// Test panic when right is nil for AND operator
func TestSelectBuilder_Where_OpAnd_NilRight_ShouldPanic(t *testing.T) {
	setupTestRegistry()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for AND with nil Right, but didn't panic")
		}
	}()

	builder := NewSelectBuilder(model.User{})
	builder.Where(&BinaryExpr{
		Operator: OpAnd,
		Left: &BinaryExpr{
			Operator: OpEq,
			Left:     &ColumnExpr{Name: "ID"},
			Right:    &LiteralExpr{Value: 1},
		},
		Right: nil, // This should panic for OpAnd
	}).Build()
}

// Test BETWEEN operator using TernaryExpr
func TestSelectBuilder_Where_OpBetween(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&TernaryExpr{
		Expr: &ColumnExpr{Name: "ID"},
		Low:  &LiteralExpr{Value: 10},
		High: &LiteralExpr{Value: 100},
	}).Build()

	if !strings.Contains(query, "WHERE id BETWEEN ? AND ?") {
		t.Errorf("Expected WHERE clause with id BETWEEN ? AND ?, got: %s", query)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}

	if args[0] != 10 || args[1] != 100 {
		t.Errorf("Expected args [10, 100], got %v", args)
	}
}

// Test NOT operator using UnaryExpr
func TestSelectBuilder_Where_OpNot(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&UnaryExpr{
		Operator: OpNot,
		Operand: &BinaryExpr{
			Operator: OpEq,
			Left:     &ColumnExpr{Name: "Name"},
			Right:    &LiteralExpr{Value: "Admin"},
		},
	}).Build()

	if !strings.Contains(query, "WHERE NOT (name = ?)") {
		t.Errorf("Expected WHERE clause with NOT (name = ?), got: %s", query)
	}

	if len(args) != 1 || args[0] != "Admin" {
		t.Errorf("Expected args ['Admin'], got %v", args)
	}
}

// Test complex query with UnaryExpr (IS NULL) combined with AND
func TestSelectBuilder_Where_IsNullWithAnd(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpAnd,
		Left: &BinaryExpr{
			Operator: OpGt,
			Left:     &ColumnExpr{Name: "ID"},
			Right:    &LiteralExpr{Value: 5},
		},
		Right: &UnaryExpr{
			Operator: OpIsNull,
			Operand:  &ColumnExpr{Name: "Email"},
		},
	}).Build()

	if !strings.Contains(query, "WHERE (id > ? AND email IS NULL)") {
		t.Errorf("Expected WHERE clause combining > and IS NULL, got: %s", query)
	}

	if len(args) != 1 || args[0] != 5 {
		t.Errorf("Expected args [5], got %v", args)
	}
}

// Test complex query with BETWEEN and OR
func TestSelectBuilder_Where_BetweenWithOr(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpOr,
		Left: &TernaryExpr{
			Expr: &ColumnExpr{Name: "ID"},
			Low:  &LiteralExpr{Value: 1},
			High: &LiteralExpr{Value: 10},
		},
		Right: &TernaryExpr{
			Expr: &ColumnExpr{Name: "ID"},
			Low:  &LiteralExpr{Value: 50},
			High: &LiteralExpr{Value: 100},
		},
	}).Build()

	if !strings.Contains(query, "WHERE (id BETWEEN ? AND ? OR id BETWEEN ? AND ?)") {
		t.Errorf("Expected WHERE clause with BETWEEN OR BETWEEN, got: %s", query)
	}

	if len(args) != 4 {
		t.Errorf("Expected 4 args, got %d", len(args))
	}

	if args[0] != 1 || args[1] != 10 || args[2] != 50 || args[3] != 100 {
		t.Errorf("Expected args [1, 10, 50, 100], got %v", args)
	}
}

// Test NOT with complex nested expression
func TestSelectBuilder_Where_NotWithComplexExpr(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	// NOT (id > 10 AND name = 'Admin')
	query, args := builder.Where(&UnaryExpr{
		Operator: OpNot,
		Operand: &BinaryExpr{
			Operator: OpAnd,
			Left: &BinaryExpr{
				Operator: OpGt,
				Left:     &ColumnExpr{Name: "ID"},
				Right:    &LiteralExpr{Value: 10},
			},
			Right: &BinaryExpr{
				Operator: OpEq,
				Left:     &ColumnExpr{Name: "Name"},
				Right:    &LiteralExpr{Value: "Admin"},
			},
		},
	}).Build()

	if !strings.Contains(query, "WHERE NOT ((id > ? AND name = ?))") {
		t.Errorf("Expected WHERE clause with NOT and complex expression, got: %s", query)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}

	if args[0] != 10 || args[1] != "Admin" {
		t.Errorf("Expected args [10, 'Admin'], got %v", args)
	}
}

// Test unsupported unary operator should panic
func TestSelectBuilder_Where_UnsupportedUnaryOperator_ShouldPanic(t *testing.T) {
	setupTestRegistry()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for unsupported unary operator, but didn't panic")
		} else {
			errMsg := fmt.Sprint(r)
			if !strings.Contains(errMsg, "unsupported unary operator") {
				t.Errorf("Expected panic message to contain 'unsupported unary operator', got: %s", errMsg)
			}
		}
	}()

	builder := NewSelectBuilder(model.User{})
	builder.Where(&UnaryExpr{
		Operator: Op("INVALID"),
		Operand:  &ColumnExpr{Name: "ID"},
	}).Build()
}

// Test IN with multiple values
func TestSelectBuilder_Where_OpIn_MultipleValues(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.Where(&BinaryExpr{
		Operator: OpIn,
		Left:     &ColumnExpr{Name: "ID"},
		Right:    &LiteralExpr{Value: []int{1, 2, 3, 4, 5}},
	}).Build()

	if !strings.Contains(query, "WHERE (id IN ?)") {
		t.Errorf("Expected WHERE clause with (id IN ?), got: %s", query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg (slice), got %d", len(args))
	}

	// Verify the slice value
	if argSlice, ok := args[0].([]int); !ok {
		t.Errorf("Expected arg to be []int, got %T", args[0])
	} else if len(argSlice) != 5 {
		t.Errorf("Expected slice with 5 elements, got %d", len(argSlice))
	}
}

// Test nested NOT expressions
func TestSelectBuilder_Where_NestedNot(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	// NOT (NOT (id > 10))  - double negation
	query, args := builder.Where(&UnaryExpr{
		Operator: OpNot,
		Operand: &UnaryExpr{
			Operator: OpNot,
			Operand: &BinaryExpr{
				Operator: OpGt,
				Left:     &ColumnExpr{Name: "ID"},
				Right:    &LiteralExpr{Value: 10},
			},
		},
	}).Build()

	if !strings.Contains(query, "WHERE NOT (NOT (id > ?))") {
		t.Errorf("Expected WHERE clause with nested NOT, got: %s", query)
	}

	if len(args) != 1 || args[0] != 10 {
		t.Errorf("Expected args [10], got %v", args)
	}
}
