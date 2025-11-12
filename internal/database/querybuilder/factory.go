package querybuilder

// BuilderFactory creates query builders
type BuilderFactory struct{}

// CreateSelect creates a SELECT query builder
func (f *BuilderFactory) CreateSelect(model any) QueryBuilder {
	return NewSelectBuilder(model)
}

// CreateInsert creates an INSERT query builder
func (f *BuilderFactory) CreateInsert(model any) QueryBuilder {
	return NewInsertBuilder(model)
}

// CreateBuilder creates a query builder based on the specified type
func (f *BuilderFactory) CreateBuilder(t SQLBuilderType, model any) QueryBuilder {
	switch t {
	case SelectType:
		return f.CreateSelect(model)
	case InsertType:
		return f.CreateInsert(model)
	default:
		return nil
	}
}
