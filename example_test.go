// Copyright 2019 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package badgerdb_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/uncle-gua/badgerdb"
)

type Item struct {
	ID       int
	Category string `badgerdbIndex:"Category"`
	Created  time.Time
}

func Example() {
	data := []Item{
		{
			ID:       0,
			Category: "blue",
			Created:  time.Now().Add(-4 * time.Hour),
		},
		{
			ID:       1,
			Category: "red",
			Created:  time.Now().Add(-3 * time.Hour),
		},
		{
			ID:       2,
			Category: "blue",
			Created:  time.Now().Add(-2 * time.Hour),
		},
		{
			ID:       3,
			Category: "blue",
			Created:  time.Now().Add(-20 * time.Minute),
		},
	}

	dir := tempdir()
	defer os.RemoveAll(dir)

	options := badgerdb.DefaultOptions
	options.Dir = dir
	options.ValueDir = dir
	store, err := badgerdb.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	// insert the data in one transaction
	err = store.Badger().Update(func(tx *badger.Txn) error {
		for i := range data {
			err := store.TxInsert(tx, data[i].ID, data[i])
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// Find all items in the blue category that have been created in the past hour
	var result []Item

	err = store.Find(&result, badgerdb.Where("Category").Eq("blue").And("Created").Ge(time.Now().Add(-1*time.Hour)))

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	fmt.Println(result[0].ID)
	// Output: 3

}
