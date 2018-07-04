// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Roberto Mier Escandon <rmescandon@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package badgercli

import (
	"encoding/json"

	"github.com/dgraph-io/badger"

	"github.com/pkg/errors"
)

const defaultDbPath = "/var/lib/badger"

type opType int

const (
	get    opType = 0
	set    opType = 1
	getAll opType = 2
)

type operation func(db *badger.DB, path []byte, obj interface{}) error

func dbExec(dbDir string, opType opType, path string, obj interface{}) error {
	if len(dbDir) == 0 {
		dbDir = defaultDbPath
	}
	opts := badger.DefaultOptions
	opts.Dir = dbDir
	opts.ValueDir = dbDir

	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	defer db.Close()

	key := []byte(path)

	for {
		var operation operation
		switch opType {
		case get:
			operation = load
		case set:
			operation = save
		case getAll:
			operation = list
		default:
			return errors.New("Invalid operation")
		}

		err := operation(db, key, obj)
		if err != nil && err == badger.ErrKeyNotFound && opType == get {
			opType = getAll
			continue
		}

		return err
	}
}

func load(db *badger.DB, key []byte, obj interface{}) error {
	return db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(key)
		if err != nil {
			return err
		}

		b, err := item.Value()
		if err != nil {
			return err
		}

		return json.Unmarshal(b, &obj)
	})
}

func save(db *badger.DB, key []byte, obj interface{}) error {
	return db.Update(func(tx *badger.Txn) error {
		b, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		return tx.Set(key, b)
	})
}

func list(db *badger.DB, key []byte, objList interface{}) error {
	return db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		iter := tx.NewIterator(opts)
		defer iter.Close()

		var items []interface{}

		for iter.Seek(key); iter.ValidForPrefix(key); iter.Next() {
			var obj interface{}
			item := iter.Item()
			b, err := item.Value()
			if err != nil {
				return err
			}

			if err := json.Unmarshal(b, &obj); err != nil {
				return err
			}

			items = append(items, obj)
		}

		objList = items
		return nil
	})
}
