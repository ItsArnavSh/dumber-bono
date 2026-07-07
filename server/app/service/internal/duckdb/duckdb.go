package duckdb

import (
	"database/sql"
	"fmt"

	_ "github.com/marcboeker/go-duckdb"
)

type DuckDB struct {
	conn *sql.DB
}

func NewDuckDB(path string) (*DuckDB, error) {
	// path can be a file (e.g. "bonobo.duckdb") or ":memory:" for in-memory only
	conn, err := sql.Open("duckdb", path)
	if err != nil {
		return nil, fmt.Errorf("open duckdb: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping duckdb: %w", err)
	}

	return &DuckDB{conn: conn}, nil
}

func (d *DuckDB) Close() error {
	return d.conn.Close()
}

func (d *DuckDB) Conn() *sql.DB {
	return d.conn
}
