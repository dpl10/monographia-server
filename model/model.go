package model

import (
	"github.com/jmoiron/sqlx"
)

type (
	// Model passes database connection pool
	Model struct {
		DB *sqlx.DB
	}
)

// NewModel creates an echo handler with access to the database connection pool
func NewModel(db *sqlx.DB) *Model {
	return &Model{
		DB: db,
	}
}
