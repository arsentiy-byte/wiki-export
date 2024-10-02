package app

import (
	"context"
	"fmt"
	"log"
	"wiki-export/internal/config"
)

type App struct {
	cfg      *config.Config
	provider *provider
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	fmt.Println(a.provider.Database())

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.migrateDatabase,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	log.Println("Loading config...")
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	a.cfg = cfg

	log.Println("Config loaded")
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	log.Println("Initializing service provider...")
	a.provider = newProvider(a.cfg)

	log.Println("Service provider initialized")
	return nil
}

func (a *App) migrateDatabase(ctx context.Context) error {
	log.Println("Migrating database...")

	if a.cfg.Storage.GetMigrationFile() == "" {
		log.Println("Skipping database migration: no migration file provided")
		return nil
	}

	if err := a.provider.Database().Migrate(ctx, a.cfg.Storage.GetMigrationFile()); err != nil {
		log.Printf("[ERROR] while migrating database: %s", err)

		return err
	}

	log.Println("Database migrated")

	return nil
}
