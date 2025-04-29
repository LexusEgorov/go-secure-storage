package storage

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	d *sql.DB
}

var db *DB = nil

func NewDB(connection string) *DB {
	if db != nil {
		return db
	}

	db = &DB{}
	db.connect(connection)

	return db
}

func (d DB) AddUser()  {}
func (d DB) GetUser()  {}
func (d DB) GetToken() {}
func (d DB) AddToken() {}

func (d *DB) connect(connection string) error {
	var err error
	d.d, err = sql.Open("pgx", connection)

	if err != nil {
		return err
	}

	err = d.d.Ping()

	if err != nil {
		return err
	}

	return nil
}
