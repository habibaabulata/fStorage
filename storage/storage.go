package storage

import (
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	DB    *leveldb.DB
	dbErr error
	once  sync.Once
)

func InitDB(path string) {
	once.Do(func() {
		DB, dbErr = leveldb.OpenFile(path, nil)
		if dbErr != nil {
			log.Fatalf("Failed to open database: %v", dbErr)
		}
	})
}

func CloseDB() {
	if DB != nil {
		dbErr = DB.Close()
		if dbErr != nil {
			log.Fatalf("Failed to close database: %v", dbErr)
		}
	}
}
