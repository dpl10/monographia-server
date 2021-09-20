package util

import (
	"time"

	// a blank import because go-lint sucks ass (and not in a good way)
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQL connects (or panics) to the MySQL database
func MySQL(x string, y int) *sqlx.DB {
	db := sqlx.MustConnect("mysql", x)
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(y)
	db.SetMaxIdleConns(y)
	return db
}
