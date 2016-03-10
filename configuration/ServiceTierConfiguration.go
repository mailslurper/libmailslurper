// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
package configuration

import (
	"github.com/adampresley/logging"
	"github.com/mailslurper/libmailslurper/middleware"
	"github.com/mailslurper/libmailslurper/storage"
)

/*
ServiceTierConfiguration allows a caller to configure how to start
and run the service tier HTTP server
*/
type ServiceTierConfiguration struct {
	Address          string
	Context          *middleware.AppContext
	Database         storage.IStorage
	Log              *logging.Logger
	Port             int
	CertFile         string
	KeyFile          string
	CertIsSelfSigned bool
}
