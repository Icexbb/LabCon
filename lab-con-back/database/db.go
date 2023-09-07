package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	dbType string
	dbPath string
}

func NewDB(dbType, dbPath string) DB {
	var db = DB{
		dbType: dbType,
		dbPath: dbPath,
	}
	return db
}
func (d DB) GetHandler() *sql.DB {
	db, _ := sql.Open(d.dbType, d.dbPath)
	return db
}
