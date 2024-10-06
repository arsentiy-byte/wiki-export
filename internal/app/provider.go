package app

import (
	"wiki-export/internal/config"
	"wiki-export/internal/repository"
	"wiki-export/internal/service"
	"wiki-export/pkg/database"
	"wiki-export/pkg/database/mysql"
)

type provider struct {
	cfg    *config.Config
	db     database.Database
	export service.ExportService
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

func (p *provider) ExportService() service.ExportService {
	if p.export == nil {
		p.export = service.NewExportService(repository.NewPageRepository(p.Database()), p.cfg.Clients.Http.Wiki)
	}

	return p.export
}
