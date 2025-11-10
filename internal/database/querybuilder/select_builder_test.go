package querybuilder

import (
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
	query, args := builder.Where("email = ?", "test@example.com").Build()

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
	query, args := builder.
		Where("email = ?", "test@example.com").
		Where("name = ?", "John").
		Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "WHERE email = ? AND name = ?") {
		t.Errorf("Unexpected query: %s", query)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
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

	expectedQuery := "SELECT id, email, name, password FROM users LIMIT 10 OFFSET 20"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}
}

func TestSelectBuilder_ComplexQuery(t *testing.T) {
	setupTestRegistry()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.
		Where("email LIKE ?", "%@example.com").
		Where("name != ?", "Admin").
		OrderBy("name", Ascending).
		OrderBy("id", Descending).
		Limit(25).
		Offset(50).
		Build()

	if !strings.Contains(query, "SELECT") || !strings.Contains(query, "FROM users") ||
		!strings.Contains(query, "WHERE email LIKE ? AND name != ?") ||
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
	result := builder.Where("id > ?", 1)
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
