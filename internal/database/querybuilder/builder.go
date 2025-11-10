package querybuilder

type QueryBuilder interface {
	Build() (string, []any)
}
