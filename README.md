# Query Builder

A Go package for building SQL queries safely and fluently using the builder pattern. This approach helps prevent SQL injection and syntax errors by abstracting raw SQL strings.

## Features

### ✅ Implemented
- **`SELECT` Query Builder**: Fluent API for `SELECT`, `WHERE`, `ORDER BY`, `LIMIT`, and `OFFSET`.
- **Expression Tree**: Supports complex, nested conditions using binary (`AND`, `OR`, `=`, `LIKE`, `IN`), unary (`NOT`, `IS NULL`), and ternary (`BETWEEN`) expressions.
- **Validation**: Automatically validates column names and transforms them from struct fields to database column names (e.g., `User.ID` -> `id`).
- **Safety**: Generates parameterized queries to prevent SQL injection.

### 🔄 Future Enhancements
- Full implementation for `INSERT`, `UPDATE`, and `DELETE` builders.
- `JOIN` support.
- `GROUP BY` and `HAVING` clauses.
- Subqueries and aggregate functions.

## Usage

### Simple SELECT Query
```go
package main

import (
    "fmt"
    "little-orm/internal/database/querybuilder"
    "little-orm/internal/model"
)

func main() {
    builder := querybuilder.NewSelectBuilder(model.User{})
    
    query, args := builder.
        Select("ID", "Name").
        Where(querybuilder.Eq("ID", 1)).
        Build()

    fmt.Printf("Query: %s\n", query)
    fmt.Printf("Args: %v\n", args)
    
    // Output:
    // Query: SELECT id, name FROM users WHERE id = ?
    // Args: [1]
}
```

### Complex SELECT Query
Build a query to find users with `ID > 10` and `Name LIKE '%John%'`, ordered by ID descending.
```go
    builder := querybuilder.NewSelectBuilder(model.User{})

    query, args := builder.
        Select("ID", "Name", "Email").
        Where(
            querybuilder.And(
                querybuilder.Gt("ID", 10),
                querybuilder.Like("Name", "%John%"),
            ),
        ).
        OrderBy("ID", querybuilder.Descending).
        Limit(5).
        Offset(10).
        Build()

    fmt.Printf("Query: %s\n", query)
    fmt.Printf("Args: %v\n", args)

    // Output:
    // Query: SELECT id, name, email FROM users WHERE (id > ? AND name LIKE ?) ORDER BY ID DESC LIMIT 5 OFFSET 10
    // Args: [10 %John%]
```

### Other Expression Examples

- **`BETWEEN`**:
  ```go
  builder.Where(querybuilder.Between("ID", 10, 100))
  // SQL: WHERE id BETWEEN ? AND ?
  ```
- **`IN`**:
  ```go
  builder.Where(querybuilder.In("ID", []int{1, 2, 3}))
  // SQL: WHERE id IN ?
  ```
- **`IS NULL`**:
  ```go
  builder.Where(querybuilder.IsNull("Email"))
  // SQL: WHERE email IS NULL
  ```
- **`NOT`**:
  ```go
  builder.Where(querybuilder.Not(querybuilder.Eq("Name", "Admin")))
  // SQL: WHERE NOT (name = ?)
  ```

## Testing

The package includes comprehensive test coverage for all operators and edge cases. To run the tests:

```bash
go test ./internal/database/querybuilder/... -v
```

## Architecture

The `querybuilder` is organized into several files, each with a distinct responsibility:

```
querybuilder/
├── builder.go          # Core QueryBuilder interface
├── expression.go       # Expression types (Binary, Unary, etc.)
├── select_builder.go   # SELECT query builder implementation
├── insert_builder.go   # (Partial) INSERT query builder
├── validate.go         # Expression validation logic
├── factory.go          # Factory for creating builders
├── const.go            # Constants for operators, types, etc.
└── helper.go           # Helper functions
```

## Design Principles

1.  **Builder Pattern**: Provides a fluent API for query construction.
2.  **Expression Tree**: Uses a recursive structure for complex, nested query conditions.
3.  **Type Safety**: Leverages Go's type system to catch errors at compile time.
4.  **Separation of Concerns**: Expression logic, query building, and validation are handled in separate components.
5.  **Validation First**: Catches invalid field names or operations early.