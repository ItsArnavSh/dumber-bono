package badger

import (
	"encoding/json"
	"strings"

	"github.com/dgraph-io/badger/v4"
)

func (b *Badger) Set(namespace string, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	fullKey := namespace + ":" + key
	return b.conn.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(fullKey), data)
	})
}

func (b *Badger) Get(namespace, key string, dest any) error {
	fullKey := namespace + ":" + key
	return b.conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(fullKey))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, dest)
		})
	})
}
func (b *Badger) List(namespace string) (map[string][]byte, error) {
	prefix := []byte(namespace + ":")
	result := make(map[string][]byte)

	err := b.conn.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			trimmedKey := strings.TrimPrefix(string(item.Key()), string(prefix))

			err := item.Value(func(val []byte) error {
				valCopy := append([]byte{}, val...) // must copy, val isn't valid outside this closure
				result[trimmedKey] = valCopy
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return result, err
}
func (b *Badger) BulkSet(namespace string, kv map[string]any) error {
	return b.conn.Update(func(txn *badger.Txn) error {
		for key, value := range kv {
			data, err := json.Marshal(value)
			if err != nil {
				return err
			}
			fullKey := namespace + ":" + key
			if err := txn.Set([]byte(fullKey), data); err != nil {
				return err
			}
		}
		return nil
	})
}
