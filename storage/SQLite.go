// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package storage

import (
	"log"

	"github.com/adampresley/golangdb"
)

func CreateSqlliteDatabase() error {
	log.Println("INFO - Creating tables...")

	var err error

	sql := `
		CREATE TABLE mailitem (
			id TEXT PRIMARY KEY,
			dateSent TEXT,
			fromAddress TEXT,
			toAddressList TEXT,
			subject TEXT,
			xmailer TEXT,
			body TEXT,
			contentType TEXT,
			boundary TEXT
		);`

	_, err = golangdb.Db["lib"].Exec(sql)
	if err != nil {
		return err
	}

	sql = `
		CREATE TABLE attachment (
			id TEXT PRIMARY KEY,
			mailItemId INTEGER,
			fileName TEXT,
			contentType TEXT,
			content TEXT
		);`

	_, err = golangdb.Db["lib"].Exec(sql)
	if err != nil {
		return err
	}

	log.Println("INFO - Created tables successfully.")
	return nil
}
