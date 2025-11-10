package querybuilder

import (
	"little-orm/internal/database/registry"
	"little-orm/internal/model"
)

// setupTestRegistry initializes the registry with test models
// This is safe to call multiple times as it just ensures the model is registered
func setupTestRegistry() {
	reg := registry.GetDBRegistry()

	// Defer recover in case model is already registered
	defer func() {
		recover()
	}()

	reg.Register(model.User{})
}

// init ensures the test model is registered when the test package loads
func init() {
	setupTestRegistry()
}
