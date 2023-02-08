package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

var maxOpenDbConn = 10
var maxIdleDbConn = 5
var maxDbLifetime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {

	d, err := NewDatabase(dsn)
	if err != nil {
		return nil, err
	}

	if err = TestDB(d); err != nil {
		return nil, err
	}

	dbConn.SQL = d
	return dbConn, nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func TestDB(d *sql.DB) error {
	if err := d.Ping(); err != nil {
		return err
	}
	return nil
}
