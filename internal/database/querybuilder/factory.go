package querybuilder

type BuilderFactory struct{}

func (f *BuilderFactory) CreateSelect(model any) QueryBuilder {
	return NewSelectBuilder(model)
}

func (f *BuilderFactory) CreateInsert(model any) QueryBuilder {
	return NewInsertBuilder(model)
}

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
