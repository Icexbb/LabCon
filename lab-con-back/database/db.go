package database

import "database/sql"

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
