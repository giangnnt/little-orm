package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-chat/internal/database/registry"
	"go-chat/internal/model"
	"strings"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

func init() {
	dbListRegistry()
}

func dbListRegistry() {
	registry := registry.GetDBRegistry()
	registry.Register(model.User{})
	registry.Register(model.Message{})
}

func GetDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open(
			"postgres",
			"postgres://postgres:123456@localhost:5432/go_chat?sslmode=disable",
		)
		if err != nil {
			panic(err)
		}
	})
	return db
}

func Get[T any](
	db *sql.DB,
	table string,
	fields []string,
	wheres []string,
	args ...any,
) (*T, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(fields, ", "),
		table,
		strings.Join(wheres, " AND "),
	)
	fmt.Println(query)
}
