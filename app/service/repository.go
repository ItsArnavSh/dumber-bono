package service

import (
	"context"
	"dubmer-bono/app/service/internal/badger"
	"dubmer-bono/app/service/internal/duckdb"
	"dubmer-bono/app/service/internal/sqlite"
	"fmt"
	"path/filepath"
)

type Repository struct {
	Cache *badger.Badger
	DB    *sqlite.SQLite
	OLAP  *duckdb.DuckDB
}

func NewRepository(ctx context.Context, root string) (*Repository, error) {
	cache, err := badger.NewBadger(filepath.Join(root, "badger.db"))
	if err != nil {
		return nil, fmt.Errorf("init badger: %w", err)
	}

	db, err := sqlite.NewSQLite(filepath.Join(root, "sqlite.db"))
	if err != nil {
		cache.Close()
		return nil, fmt.Errorf("init sqlite: %w", err)
	}

	olap, err := duckdb.NewDuckDB(ctx, filepath.Join(root, "duck.db"))
	if err != nil {
		cache.Close()
		db.Close()
		return nil, fmt.Errorf("init duckdb: %w", err)
	}

	return &Repository{
		Cache: cache,
		DB:    db,
		OLAP:  olap,
	}, nil
}

func (r *Repository) Close() error {
	var errs []error

	if err := r.Cache.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close badger: %w", err))
	}
	if err := r.DB.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close sqlite: %w", err))
	}
	if err := r.OLAP.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close duckdb: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("repository close errors: %v", errs)
	}
	return nil
}
