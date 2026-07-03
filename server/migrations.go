package main

import (
	"database/sql"
	"embed"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose/v3"
)

//go:embed infra/clickhouse/*.sql
var embedMigrations embed.FS

func RunClickHouseMigrations(dsn string) error {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("clickhouse"); err != nil {
		return err
	}

	return goose.Up(db, "infra/clickhouse")
}
