package app

import (
	"wiki-export/internal/config"
	"wiki-export/pkg/database"
	"wiki-export/pkg/database/mysql"
)

type provider struct {
	cfg *config.Config
	db  database.Database
}

func newProvider(cfg *config.Config) *provider {
	return &provider{
		cfg: cfg,
	}
}

func (p *provider) Database() database.Database {
	if p.db == nil {
		var err error
		db, err := mysql.NewDatabase(p.cfg.Storage)
		if err != nil {
			panic(err)
		}

		p.db = db
	}

	return p.db
}
