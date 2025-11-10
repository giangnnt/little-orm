package registry

import (
	"reflect"
	"sync"
	"testing"
)

// TestModel for testing purposes
type TestModel struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

// AnotherTestModel for testing multiple models
type AnotherTestModel struct {
	ID    int    `db:"id"`
	Title string `db:"title"`
}

// ModelWithoutTags for testing models without db tags
type ModelWithoutTags struct {
	ID   int
	Name string
}

// ModelWithPartialTags for testing models with some fields having db tags
type ModelWithPartialTags struct {
	ID   int    `db:"id"`
	Name string // no db tag
	Age  int    `db:"age"`
}

// resetRegistry resets the singleton instance for testing
func resetRegistry() {
	instance = nil
	once = sync.Once{}
}

func TestGetDBRegistry_Singleton(t *testing.T) {
	resetRegistry()

	reg1 := GetDBRegistry()
	reg2 := GetDBRegistry()

	if reg1 != reg2 {
		t.Error("Expected GetDBRegistry to return the same instance (singleton)")
	}

	if reg1 == nil {
		t.Fatal("Expected registry to not be nil")
	}

	if reg1.cache == nil {
		t.Error("Expected cache to be initialized")
	}
}

func TestGetDBRegistry_Initialization(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()

	if reg.cache == nil {
		t.Error("Expected cache to be initialized")
	}

	if len(reg.cache) != 0 {
		t.Error("Expected cache to be empty on initialization")
	}
}

func TestDBRegistry_Register(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()
	reg.Register(TestModel{})

	if len(reg.cache) != 1 {
		t.Errorf("Expected 1 table in cache, got %d", len(reg.cache))
	}

	tableMeta, ok := reg.cache["testmodels"]
	if !ok {
		t.Fatal("Expected testmodels to be registered in cache")
	}

	if tableMeta.TableName != "testmodels" {
		t.Errorf("Expected table name 'testmodels', got '%s'", tableMeta.TableName)
	}

	expectedColumns := 4
	if len(tableMeta.Columns) != expectedColumns {
		t.Errorf("Expected %d columns, got %d", expectedColumns, len(tableMeta.Columns))
	}

	// Check specific columns
	expectedCols := map[string]string{
		"ID":    "id",
		"Name":  "name",
		"Email": "email",
		"Age":   "age",
	}

	for fieldName, dbTag := range expectedCols {
		col, ok := tableMeta.Columns[fieldName]
		if !ok {
			t.Errorf("Expected column '%s' to be registered", fieldName)
			continue
		}

		if col.DBTag != dbTag {
			t.Errorf("Expected column '%s' to have DBTag '%s', got '%s'", fieldName, dbTag, col.DBTag)
		}

		if col.Name != fieldName {
			t.Errorf("Expected column Name to be '%s', got '%s'", fieldName, col.Name)
		}
	}
}

func TestDBRegistry_Register_Pointer(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()
	reg.Register(&TestModel{})

	tableMeta, ok := reg.cache["testmodels"]
	if !ok {
		t.Fatal("Expected testmodels to be registered when passed as pointer")
	}

	if tableMeta.TableName != "testmodels" {
		t.Errorf("Expected table name 'testmodels', got '%s'", tableMeta.TableName)
	}
}

func TestDBRegistry_Register_MultipleModels(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()
	reg.Register(TestModel{})
	reg.Register(AnotherTestModel{})

	if len(reg.cache) != 2 {
		t.Errorf("Expected 2 tables in cache, got %d", len(reg.cache))
	}

	_, ok1 := reg.cache["testmodels"]
	if !ok1 {
		t.Error("Expected testmodels to be registered")
	}

	_, ok2 := reg.cache["anothertestmodels"]
	if !ok2 {
		t.Error("Expected anothertestmodels to be registered")
	}
}

func TestDBRegistry_Register_WithoutTags(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()
	reg.Register(ModelWithoutTags{})

	tableMeta, ok := reg.cache["modelwithouttagss"]
	if !ok {
		t.Fatal("Expected modelwithouttagss to be registered")
	}

	// Should have no columns since no db tags
	if len(tableMeta.Columns) != 0 {
		t.Errorf("Expected 0 columns (no db tags), got %d", len(tableMeta.Columns))
	}
}

func TestDBRegistry_Register_PartialTags(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()
	reg.Register(ModelWithPartialTags{})

	tableMeta, ok := reg.cache["modelwithpartialtagss"]
	if !ok {
		t.Fatal("Expected modelwithpartialtagss to be registered")
	}

	// Should only have columns with db tags
	if len(tableMeta.Columns) != 2 {
		t.Errorf("Expected 2 columns (only those with db tags), got %d", len(tableMeta.Columns))
	}

	_, hasID := tableMeta.Columns["ID"]
	if !hasID {
		t.Error("Expected ID column to be registered")
	}

	_, hasAge := tableMeta.Columns["Age"]
	if !hasAge {
		t.Error("Expected Age column to be registered")
	}

	_, hasName := tableMeta.Columns["Name"]
	if hasName {
		t.Error("Expected Name column NOT to be registered (no db tag)")
	}
}

func TestDBRegistry_GetTableMeta(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()
	reg.Register(TestModel{})

	tableMeta := reg.GetTableMeta(TestModel{})

	if tableMeta.TableName != "testmodels" {
		t.Errorf("Expected table name 'testmodels', got '%s'", tableMeta.TableName)
	}

	if len(tableMeta.Columns) != 4 {
		t.Errorf("Expected 4 columns, got %d", len(tableMeta.Columns))
	}
}

func TestDBRegistry_GetTableMeta_Pointer(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()
	reg.Register(TestModel{})

	tableMeta := reg.GetTableMeta(&TestModel{})

	if tableMeta.TableName != "testmodels" {
		t.Errorf("Expected table name 'testmodels', got '%s'", tableMeta.TableName)
	}
}

func TestDBRegistry_GetTableMeta_NotRegistered(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for unregistered model, but didn't panic")
		}
	}()

	reg.GetTableMeta(TestModel{})
}

func TestDBRegistry_ThreadSafety(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()
	var wg sync.WaitGroup

	// Register multiple models concurrently
	models := []interface{}{
		TestModel{},
		AnotherTestModel{},
		ModelWithPartialTags{},
	}

	for _, model := range models {
		wg.Add(1)
		go func(m interface{}) {
			defer wg.Done()
			reg.Register(m)
		}(model)
	}

	wg.Wait()

	// Verify all models are registered
	if len(reg.cache) != 3 {
		t.Errorf("Expected 3 models registered, got %d", len(reg.cache))
	}

	// Test concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reg.GetTableMeta(TestModel{})
		}()
	}

	wg.Wait()
}

func TestGetTableName(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		expected string
	}{
		{
			name:     "Simple model",
			typeName: "User",
			expected: "users",
		},
		{
			name:     "CamelCase model",
			typeName: "TestModel",
			expected: "testmodels",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy type with the name
			var model interface{}
			if tt.typeName == "User" {
				type User struct{}
				model = User{}
			} else if tt.typeName == "TestModel" {
				model = TestModel{}
			}

			typ := reflect.TypeOf(model)
			result := getTableName(&typ)

			if result != tt.expected {
				t.Errorf("Expected table name '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetTableColsNameMap(t *testing.T) {
	typ := reflect.TypeOf(TestModel{})
	cols := getTableColsNameMap(&typ)

	if len(cols) != 4 {
		t.Errorf("Expected 4 columns, got %d", len(cols))
	}

	// Test ID column
	idCol, ok := cols["ID"]
	if !ok {
		t.Fatal("Expected ID column to exist")
	}

	if idCol.DBTag != "id" {
		t.Errorf("Expected DBTag 'id', got '%s'", idCol.DBTag)
	}

	if idCol.Name != "ID" {
		t.Errorf("Expected Name 'ID', got '%s'", idCol.Name)
	}

	if idCol.Type != "int" {
		t.Errorf("Expected Type 'int', got '%s'", idCol.Type)
	}

	if idCol.Tag == "" {
		t.Error("Expected Tag to be non-empty")
	}

	// Test Name column
	nameCol, ok := cols["Name"]
	if !ok {
		t.Fatal("Expected Name column to exist")
	}

	if nameCol.DBTag != "name" {
		t.Errorf("Expected DBTag 'name', got '%s'", nameCol.DBTag)
	}

	if nameCol.Type != "string" {
		t.Errorf("Expected Type 'string', got '%s'", nameCol.Type)
	}
}

func TestGetTableColsNameMap_IgnoreFieldsWithoutTags(t *testing.T) {
	typ := reflect.TypeOf(ModelWithPartialTags{})
	cols := getTableColsNameMap(&typ)

	// Should only have 2 columns (ID and Age with db tags)
	if len(cols) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(cols))
	}

	_, hasID := cols["ID"]
	if !hasID {
		t.Error("Expected ID column to exist")
	}

	_, hasAge := cols["Age"]
	if !hasAge {
		t.Error("Expected Age column to exist")
	}

	_, hasName := cols["Name"]
	if hasName {
		t.Error("Expected Name column NOT to exist (no db tag)")
	}
}

func TestColumnMeta_Fields(t *testing.T) {
	typ := reflect.TypeOf(TestModel{})
	cols := getTableColsNameMap(&typ)

	emailCol := cols["Email"]

	// Test all fields of ColumnMeta
	if emailCol.DBTag != "email" {
		t.Errorf("Expected DBTag 'email', got '%s'", emailCol.DBTag)
	}

	if emailCol.Name != "Email" {
		t.Errorf("Expected Name 'Email', got '%s'", emailCol.Name)
	}

	if emailCol.Type != "string" {
		t.Errorf("Expected Type 'string', got '%s'", emailCol.Type)
	}

	// Tag should contain the full struct tag
	if emailCol.Tag == "" {
		t.Error("Expected Tag to be non-empty")
	}

	// Tag should contain db:"email"
	if !contains(emailCol.Tag, `db:"email"`) {
		t.Errorf("Expected Tag to contain 'db:\"email\"', got '%s'", emailCol.Tag)
	}
}

func TestTableMeta_Fields(t *testing.T) {
	resetRegistry()

	reg := GetDBRegistry()
	reg.Register(TestModel{})

	tableMeta := reg.GetTableMeta(TestModel{})

	// Test TableName
	if tableMeta.TableName != "testmodels" {
		t.Errorf("Expected TableName 'testmodels', got '%s'", tableMeta.TableName)
	}

	// Test Columns map exists
	if tableMeta.Columns == nil {
		t.Fatal("Expected Columns to be non-nil")
	}

	// Test Columns map has correct entries
	if len(tableMeta.Columns) != 4 {
		t.Errorf("Expected 4 columns, got %d", len(tableMeta.Columns))
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
