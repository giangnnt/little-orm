package querybuilder

// QueryBuilder interface for building SQL queries
type QueryBuilder interface {
	// Build constructs the final SQL query and returns it with arguments
	Build() (string, []any)
}
