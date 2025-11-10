package querybuilder

type BuilderFactory struct{}

func (f *BuilderFactory) Create(t SQLBuilderType, model any) QueryBuilder {
	switch t {
	case SelectType:
		return NewSelectBuilder(model)
	// case InsertType:
	// 	return &InsertBuilder{}
	default:
		return nil
	}
}
