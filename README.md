# Query Builder Package

Package `querybuilder` cung cấp các công cụ để xây dựng SQL queries một cách an toàn và linh hoạt thông qua builder pattern.

## Tổng quan

Query Builder giúp xây dựng SQL queries mà không cần viết raw SQL strings, giảm thiểu SQL injection và lỗi cú pháp.

## Components

### 1. Expression System

#### Expr Interface
Interface cơ bản cho tất cả các expressions:
```go
type Expr interface {
    ToSQL() (string, []any)
}
```

#### Expression Types

**ColumnExpr** - Đại diện cho column reference
```go
&ColumnExpr{Name: "ID"}
// SQL: id (sau khi transform)
```

**LiteralExpr** - Đại diện cho literal value
```go
&LiteralExpr{Value: 42}
// SQL: ? với args: [42]
```

**BinaryExpr** - Đại diện cho binary operations (=, !=, >, <, AND, OR, IN, etc.)
```go
&BinaryExpr{
    Operator: OpEq,
    Left:     &ColumnExpr{Name: "ID"},
    Right:    &LiteralExpr{Value: 1},
}
// SQL: id = ? với args: [1]
```

**UnaryExpr** - Đại diện cho unary operations (IS NULL, IS NOT NULL, NOT)
```go
&UnaryExpr{
    Operator: "IS NULL",
    Operand:  &ColumnExpr{Name: "Email"},
}
// SQL: email IS NULL
```

**TernaryExpr** - Đại diện cho ternary operations (BETWEEN)
```go
&TernaryExpr{
    Expr: &ColumnExpr{Name: "ID"},
    Low:  &LiteralExpr{Value: 10},
    High: &LiteralExpr{Value: 100},
}
// SQL: id BETWEEN ? AND ? với args: [10, 100]
```

### 2. Query Builders

#### SelectBuilder

Xây dựng SELECT queries với fluent API.

**Khởi tạo:**
```go
builder := NewSelectBuilder(model.User{})
```

**API Methods:**

- `Select(fields ...string)` - Chọn các fields cụ thể (mặc định: tất cả fields)
- `Where(expr Expr)` - Thêm WHERE clause
- `OrderBy(field string, order SortOrder)` - Thêm ORDER BY clause
- `Limit(n int)` - Thêm LIMIT clause
- `Offset(m int)` - Thêm OFFSET clause
- `Build()` - Tạo SQL query cuối cùng

**Ví dụ đơn giản:**
```go
builder := NewSelectBuilder(model.User{})
query, args := builder.
    Select("ID", "Name").
    Where(&BinaryExpr{
        Operator: OpEq,
        Left:     &ColumnExpr{Name: "ID"},
        Right:    &LiteralExpr{Value: 1},
    }).
    Build()
// SQL: SELECT id, name FROM users WHERE id = ?
// Args: [1]
```

**Ví dụ phức tạp:**
```go
builder := NewSelectBuilder(model.User{})
query, args := builder.
    Where(&BinaryExpr{
        Operator: OpAnd,
        Left: &BinaryExpr{
            Operator: OpGt,
            Left:     &ColumnExpr{Name: "ID"},
            Right:    &LiteralExpr{Value: 10},
        },
        Right: &BinaryExpr{
            Operator: OpLike,
            Left:     &ColumnExpr{Name: "Name"},
            Right:    &LiteralExpr{Value: "%John%"},
        },
    }).
    OrderBy("ID", Ascending).
    Limit(10).
    Offset(5).
    Build()
// SQL: SELECT id, email, name, password FROM users
//      WHERE (id > ? AND name LIKE ?)
//      ORDER BY ID ASC LIMIT 10 OFFSET 5
// Args: [10, "%John%"]
```

**BETWEEN example:**
```go
builder := NewSelectBuilder(model.User{})
query, args := builder.
    Where(&TernaryExpr{
        Expr: &ColumnExpr{Name: "ID"},
        Low:  &LiteralExpr{Value: 10},
        High: &LiteralExpr{Value: 100},
    }).
    Build()
// SQL: SELECT id, email, name, password FROM users WHERE id BETWEEN ? AND ?
// Args: [10, 100]
```

**IS NULL example:**
```go
builder := NewSelectBuilder(model.User{})
query, args := builder.
    Where(&UnaryExpr{
        Operator: "IS NULL",
        Operand:  &ColumnExpr{Name: "Email"},
    }).
    Build()
// SQL: SELECT id, email, name, password FROM users WHERE email IS NULL
```

**NOT example:**
```go
builder := NewSelectBuilder(model.User{})
query, args := builder.
    Where(&UnaryExpr{
        Operator: OpNot,
        Operand: &BinaryExpr{
            Operator: OpEq,
            Left:     &ColumnExpr{Name: "Name"},
            Right:    &LiteralExpr{Value: "Admin"},
        },
    }).
    Build()
// SQL: SELECT id, email, name, password FROM users WHERE NOT (name = ?)
// Args: ["Admin"]
```

**IN example:**
```go
builder := NewSelectBuilder(model.User{})
query, args := builder.
    Where(&BinaryExpr{
        Operator: OpIn,
        Left:     &ColumnExpr{Name: "ID"},
        Right:    &LiteralExpr{Value: []int{1, 2, 3}},
    }).
    Build()
// SQL: SELECT id, email, name, password FROM users WHERE (id IN ?)
// Args: [[1, 2, 3]]
```

#### InsertBuilder

Xây dựng INSERT queries (TODO: chưa implement đầy đủ).

```go
builder := NewInsertBuilder(model.User{})
query, args := builder.Build()
```

### 3. Validation & Transformation

#### ExprValidator

Validates expressions và transforms column names sang database column names.

**Chức năng:**
- Kiểm tra column có tồn tại trong table không
- Transform column name từ struct field → database column tag
- Validate đệ quy toàn bộ expression tree

**Ví dụ:**
```go
// Input: ColumnExpr{Name: "ID"}
// Output: ColumnExpr{Name: "id"} (after validation)
```

### 4. Factory Pattern

#### BuilderFactory

Factory để tạo query builders.

```go
factory := &BuilderFactory{}

// Tạo SelectBuilder
selectBuilder := factory.CreateSelect(model.User{})

// Tạo InsertBuilder
insertBuilder := factory.CreateInsert(model.User{})

// Tạo builder theo type
builder := factory.CreateBuilder(SelectType, model.User{})
```

### 5. Constants & Types

#### Operators
```go
// Comparison operators
OpEq      = "="
OpNEq     = "!="
OpGt      = ">"
OpLt      = "<"
OpGte     = ">="
OpLte     = "<="
OpLike    = "LIKE"
OpIn      = "IN"
OpNIn     = "NOT IN"
OpIsNull  = "IS NULL"
OpIsNNull = "IS NOT NULL"
OpBetween = "BETWEEN"

// Logical operators
OpAnd = "AND"
OpOr  = "OR"
OpNot = "NOT"
```

#### Sort Orders
```go
Ascending  = "ASC"
Descending = "DESC"
```

#### Builder Types
```go
SelectType = "select"
InsertType = "insert"
```

## Features

### ✅ Implemented

1. **SELECT Query Builder**
   - Basic SELECT with all fields
   - SELECT với specific fields
   - WHERE clauses với complex expressions
   - ORDER BY (single và multiple columns)
   - LIMIT và OFFSET
   - Expression system (Binary, Unary, Ternary)

2. **Expression System**
   - BinaryExpr: =, !=, >, <, >=, <=, LIKE, IN, NOT IN, AND, OR
   - UnaryExpr: IS NULL, IS NOT NULL, NOT
   - TernaryExpr: BETWEEN
   - Nested expressions với unlimited depth

3. **Validation**
   - Column existence validation
   - Column name transformation (struct field → db tag)
   - Expression tree validation

4. **Safety Features**
   - SQL injection protection (parameterized queries)
   - Panic on invalid operations
   - Type safety

### 🔄 Future Enhancements

1. **INSERT Builder** - Full implementation
2. **UPDATE Builder** - Chưa implement
3. **DELETE Builder** - Chưa implement
4. **JOIN Support** - Chưa support
5. **GROUP BY và HAVING** - Chưa implement
6. **Subqueries** - Chưa support
7. **Aggregate Functions** - Chưa support

## Edge Cases & Behaviors

### Documented Behaviors

1. **Multiple Where() calls**: Last call overwrites previous
2. **Negative Limit/Offset**: Silently ignored (không thêm vào query)
3. **Zero Limit/Offset**: Không xuất hiện trong query
4. **Nil operands**: Panic (defensive programming)
5. **Unsupported operators**: BinaryExpr returns empty string
6. **Invalid field names**:
   - Select(): Panic
   - OrderBy(): No validation (pass through)

## Testing

Package có comprehensive test coverage:
- 62 test cases
- Coverage cho tất cả operators
- Edge cases testing
- Panic recovery tests

Chạy tests:
```bash
go test ./internal/database/querybuilder/... -v
```

## Usage Example - Complete Workflow

```go
package main

import (
    "fmt"
    "little-orm/internal/database/querybuilder"
    "little-orm/internal/model"
)

func main() {
    // Khởi tạo builder
    builder := querybuilder.NewSelectBuilder(model.User{})

    // Complex query: Tìm users có ID > 10 và name LIKE '%John%'
    // Order by ID descending, limit 5, offset 10
    query, args := builder.
        Select("ID", "Name", "Email").
        Where(&querybuilder.BinaryExpr{
            Operator: querybuilder.OpAnd,
            Left: &querybuilder.BinaryExpr{
                Operator: querybuilder.OpGt,
                Left:     &querybuilder.ColumnExpr{Name: "ID"},
                Right:    &querybuilder.LiteralExpr{Value: 10},
            },
            Right: &querybuilder.BinaryExpr{
                Operator: querybuilder.OpLike,
                Left:     &querybuilder.ColumnExpr{Name: "Name"},
                Right:    &querybuilder.LiteralExpr{Value: "%John%"},
            },
        }).
        OrderBy("ID", querybuilder.Descending).
        Limit(5).
        Offset(10).
        Build()

    fmt.Printf("Query: %s\n", query)
    fmt.Printf("Args: %v\n", args)

    // Output:
    // Query: SELECT id, name, email FROM users WHERE (id > ? AND name LIKE ?) ORDER BY ID DESC LIMIT 5 OFFSET 10
    // Args: [10 %John%]
}
```

## Architecture

```
querybuilder/
├── builder.go          # QueryBuilder interface
├── expression.go       # Expression types (Binary, Unary, Ternary, Column, Literal)
├── select_builder.go   # SELECT query builder
├── insert_builder.go   # INSERT query builder (partial)
├── validate.go         # Expression validator
├── factory.go          # Builder factory
├── const.go           # Constants (operators, types)
└── helper.go          # Helper functions
```

## Design Principles

1. **Builder Pattern**: Fluent API cho query construction
2. **Expression Tree**: Recursive structure cho complex queries
3. **Type Safety**: Strong typing với Go types
4. **Immutability**: Builder methods return new builder (chainable)
5. **Validation**: Early validation với meaningful errors
6. **Separation of Concerns**: Expression, Building, Validation tách biệt

## Contributing

Khi thêm features mới:
1. Implement expression type nếu cần
2. Add validation logic
3. Viết comprehensive tests (including edge cases)
4. Update README
5. Add godoc comments

## License

Internal package cho little-orm project.
