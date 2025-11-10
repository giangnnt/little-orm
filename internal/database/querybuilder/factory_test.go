package querybuilder

import (
	"little-orm/internal/model"
	"testing"
)

func TestBuilderFactory_CreateSelect(t *testing.T) {
	setupTestRegistry()

	factory := &BuilderFactory{}
	builder := factory.CreateSelect(model.User{})

	if builder == nil {
		t.Fatal("Expected builder to not be nil")
	}

	// Verify it returns a QueryBuilder interface
	_, ok := builder.(QueryBuilder)
	if !ok {
		t.Error("Expected builder to implement QueryBuilder interface")
	}

	// Verify it's actually a SelectBuilder
	selectBuilder, ok := builder.(*SelectBuilder)
	if !ok {
		t.Error("Expected builder to be *SelectBuilder")
	}

	if selectBuilder.table != "users" {
		t.Errorf("Expected table to be 'users', got '%s'", selectBuilder.table)
	}
}

func TestBuilderFactory_CreateInsert(t *testing.T) {
	setupTestRegistry()

	factory := &BuilderFactory{}
	builder := factory.CreateInsert(model.User{})

	if builder == nil {
		t.Fatal("Expected builder to not be nil")
	}

	// Verify it returns a QueryBuilder interface
	_, ok := builder.(QueryBuilder)
	if !ok {
		t.Error("Expected builder to implement QueryBuilder interface")
	}

	// Verify it's actually an InsertBuilder
	insertBuilder, ok := builder.(*InsertBuilder)
	if !ok {
		t.Error("Expected builder to be *InsertBuilder")
	}

	if insertBuilder.table != "users" {
		t.Errorf("Expected table to be 'users', got '%s'", insertBuilder.table)
	}
}

func TestBuilderFactory_Create_SelectType(t *testing.T) {
	setupTestRegistry()

	factory := &BuilderFactory{}
	builder := factory.Create(SelectType, model.User{})

	if builder == nil {
		t.Fatal("Expected builder to not be nil")
	}

	_, ok := builder.(*SelectBuilder)
	if !ok {
		t.Error("Expected builder to be *SelectBuilder when type is SelectType")
	}
}

func TestBuilderFactory_Create_InsertType(t *testing.T) {
	setupTestRegistry()

	factory := &BuilderFactory{}
	builder := factory.Create(InsertType, model.User{})

	if builder == nil {
		t.Fatal("Expected builder to not be nil")
	}

	_, ok := builder.(*InsertBuilder)
	if !ok {
		t.Error("Expected builder to be *InsertBuilder when type is InsertType")
	}
}

func TestBuilderFactory_Create_InvalidType(t *testing.T) {
	setupTestRegistry()

	factory := &BuilderFactory{}
	builder := factory.Create(SQLBuilderType("invalid"), model.User{})

	if builder != nil {
		t.Error("Expected nil builder for invalid type")
	}
}

func TestBuilderFactory_Create_AllTypes(t *testing.T) {
	setupTestRegistry()

	factory := &BuilderFactory{}

	testCases := []struct {
		name        string
		builderType SQLBuilderType
		expectNil   bool
		expectType  string
	}{
		{
			name:        "Select type",
			builderType: SelectType,
			expectNil:   false,
			expectType:  "*querybuilder.SelectBuilder",
		},
		{
			name:        "Insert type",
			builderType: InsertType,
			expectNil:   false,
			expectType:  "*querybuilder.InsertBuilder",
		},
		{
			name:        "Unknown type",
			builderType: SQLBuilderType("unknown"),
			expectNil:   true,
			expectType:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := factory.Create(tc.builderType, model.User{})

			if tc.expectNil {
				if builder != nil {
					t.Errorf("Expected nil builder, got %T", builder)
				}
			} else {
				if builder == nil {
					t.Fatal("Expected non-nil builder")
				}

				_, canBuild := builder.(QueryBuilder)
				if !canBuild {
					t.Error("Expected builder to implement QueryBuilder interface")
				}
			}
		})
	}
}
