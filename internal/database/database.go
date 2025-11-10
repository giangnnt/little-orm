package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"little-orm/internal/database/registry"
	"little-orm/internal/model"
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
