package registry

import (
	"reflect"
	"strings"
	"sync"
)

var (
	instance *DBRegistry
	once     sync.Once
)

type DBRegistry struct {
	mu    sync.RWMutex
	cache map[string]TableMeta
}

func GetDBRegistry() *DBRegistry {
	once.Do(func() {
		instance = &DBRegistry{cache: make(map[string]TableMeta)}
	})
	return instance
}

func (r *DBRegistry) Register(model any) {
	// get table type
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	tableName := getTableName(&t)
	tableCols := getTableCols(&t)

	meta := TableMeta{TableName: tableName, Columns: tableCols}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[t.Name()] = meta
}

func getTableName(t *reflect.Type) string {
	return strings.ToLower((*t).Name()) + "s"
}

func getTableCols(t *reflect.Type) []ColumnMeta {
	var cols []ColumnMeta
	for i := 0; i < (*t).NumField(); i++ {
		f := (*t).Field(i)
		dbTag := f.Tag.Get("db")
		if dbTag == "" {
			continue
		}
		cols = append(cols, ColumnMeta{
			Name: dbTag,
			Type: f.Type.Name(),
			Tag:  string(f.Tag),
		})
	}
	return cols
}
