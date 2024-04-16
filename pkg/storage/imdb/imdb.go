package imdb

import (
	"go.uber.org/zap"
	lg "main/pkg/logger"
	"main/resource"
	"sync"
)

type InMemoryDB struct {
	mu   sync.Mutex
	Urls map[string]string
}

var MemoryDB *InMemoryDB

func init() {
	if resource.CFG.DB == "inmemory" {
		MemoryDB = NewInMemoryDB()
	}
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Urls: make(map[string]string),
	}
}

func (db *InMemoryDB) Set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Urls[key] = value
	lg.Logger.Info("data has been saved successfully", zap.String("key", key), zap.String("value", value))

}

func (db *InMemoryDB) Get(key string) (string, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	value, ok := db.Urls[key]
	return value, ok
}
