package querybuilder

import (
	"little-orm/internal/model"
	"testing"
)

func TestNewInsertBuilder(t *testing.T) {
	setupTestRegistry()

	builder := NewInsertBuilder(model.User{})

	if builder == nil {
		t.Fatal("Expected builder to not be nil")
	}

	if builder.table != "users" {
		t.Errorf("Expected table to be 'users', got '%s'", builder.table)
	}

	if builder.columns == nil {
		t.Error("Expected columns to be initialized")
	}

	if builder.values == nil {
		t.Error("Expected values to be initialized")
	}
}

func TestInsertBuilder_Build(t *testing.T) {
	setupTestRegistry()

	builder := NewInsertBuilder(model.User{})
	query, args := builder.Build()

	// Since Build() is not implemented yet, it should return empty string and nil
	if query != "" {
		t.Errorf("Expected empty query (not implemented), got: %s", query)
	}

	if args != nil {
		t.Errorf("Expected nil args (not implemented), got: %v", args)
	}
}

// Test for future implementation
func TestInsertBuilder_TableMeta(t *testing.T) {
	setupTestRegistry()

	builder := NewInsertBuilder(model.User{})

	if builder.tableMeta.TableName != "users" {
		t.Errorf("Expected tableMeta.TableName to be 'users', got '%s'", builder.tableMeta.TableName)
	}

	if len(builder.tableMeta.Columns) == 0 {
		t.Error("Expected tableMeta.Columns to be populated")
	}

	// Check that columns are registered
	expectedColumns := []string{"ID", "Email", "Name", "Password"}
	for _, colName := range expectedColumns {
		if _, ok := builder.tableMeta.Columns[colName]; !ok {
			t.Errorf("Expected column '%s' to be registered", colName)
		}
	}
}
