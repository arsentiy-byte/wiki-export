package database

import (
	"context"
	"database/sql"
)

type Database interface {
	GetInstance() *sql.DB
	Migrate(ctx context.Context, filepath string) error
}
