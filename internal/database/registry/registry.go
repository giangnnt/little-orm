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

func (r *DBRegistry) GetTableMeta(model any) TableMeta {
	// get table type
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	tableName := getTableName(&t)

	tableMeta, ok := r.cache[tableName]
	if !ok {
		panic("Model not registered in DBRegistry")
	}
	return tableMeta

}

func (r *DBRegistry) Register(model any) {
	// get table type
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	tableName := getTableName(&t)
	tableCols := getTableColsNameMap(&t)

	tableMeta := TableMeta{TableName: tableName, Columns: tableCols}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[t.Name()] = tableMeta
}

func getTableName(t *reflect.Type) string {
	return strings.ToLower((*t).Name()) + "s"
}

// return map off cols name
func getTableColsNameMap(t *reflect.Type) map[string]ColumnMeta {
	colsMap := make(map[string]ColumnMeta, (*t).NumField())
	for i := 0; i < (*t).NumField(); i++ {
		f := (*t).Field(i)
		dbTag := f.Tag.Get("db")
		if dbTag == "" {
			continue
		}
		colsMap[f.Name] = ColumnMeta{
			DBTag: dbTag,
			Name:  f.Name,
			Type:  f.Type.String(),
			Tag:   string(f.Tag),
		}
	}
	return colsMap
}
