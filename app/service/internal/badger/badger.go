package badger

import (
	bdg "github.com/dgraph-io/badger/v4"
)

type Badger struct {
	conn *bdg.DB
}

func NewBadger(path string) (*Badger, error) {
	conn, err := bdg.Open(bdg.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return &Badger{conn: conn}, nil
}

func (b *Badger) Close() error {
	return b.conn.Close()
}
