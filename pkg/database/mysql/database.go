package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"wiki-export/pkg/config"
	"wiki-export/pkg/database"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDatabase struct {
	db *sql.DB
}

var (
	once       sync.Once
	onceError  error
	dbInstance *mysqlDatabase
)

func NewDatabase(cfg config.Storage) (database.Database, error) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&maxAllowedPacket=1073741824&multiStatements=true",
			cfg.GetUser(),
			cfg.GetPassword(),
			cfg.GetHost(),
			cfg.GetPort(),
			cfg.GetDatabase(),
		)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			onceError = err
			return
		}

		if err = db.Ping(); err != nil {
			onceError = err
			return
		}

		dbInstance = &mysqlDatabase{db: db}
	})

	if onceError != nil {
		return nil, onceError
	}

	return dbInstance, nil
}

func (m *mysqlDatabase) GetInstance() *sql.DB {
	return m.db
}

func (m *mysqlDatabase) Migrate(ctx context.Context, filepath string) error {
	dump, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = m.db.ExecContext(ctx, string(dump))
	if err != nil {
		return err
	}

	return nil
}
