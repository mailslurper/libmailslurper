// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package storage

import "log"

/*
ConnectToStorage establishes a connection to the configured database engine and returns
an object.
*/
func ConnectToStorage(storageType StorageType, connectionInfo *ConnectionInformation) (IStorage, error) {
	var err error
	var storageHandle IStorage

	log.Println("libmailslurper: INFO - Connecting to database")

	switch storageType {
	case STORAGE_SQLITE:
		storageHandle = NewSQLiteStorage(connectionInfo)

	case STORAGE_MSSQL:
		storageHandle = NewMSSQLStorage(connectionInfo)

	case STORAGE_MYSQL:
		storageHandle = NewMySQLStorage(connectionInfo)
	}

	if err = storageHandle.Connect(); err != nil {
		return storageHandle, err
	}

	if err = storageHandle.Create(); err != nil {
		return storageHandle, err
	}

	return storageHandle, nil
}
