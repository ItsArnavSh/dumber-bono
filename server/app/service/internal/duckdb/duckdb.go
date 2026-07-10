package duckdb

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/marcboeker/go-duckdb"
)

type DuckDB struct {
	conn  *sql.DB
	batch BatchedFrames
}

func NewDuckDB(ctx context.Context, path string) (*DuckDB, error) {
	// path can be a file (e.g. "bonobo.duckdb") or ":memory:" for in-memory only
	conn, err := sql.Open("duckdb", path)
	if err != nil {
		return nil, fmt.Errorf("open duckdb: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping duckdb: %w", err)
	}

	ddb := &DuckDB{conn: conn, batch: BatchedFrames{}}
	err = ddb.ApplyMigrations(ctx)
	fmt.Println("Error: ", err)
	go ddb.BatchProcess(ctx)
	return ddb, nil
}

func (d *DuckDB) Close() error {
	return d.conn.Close()
}

func (d *DuckDB) Conn() *sql.DB {
	return d.conn
}
