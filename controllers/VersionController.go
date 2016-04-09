// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/adampresley/logging"
	"github.com/gorilla/context"
	"github.com/mailslurper/libmailslurper/server"
)

/*
GetVersion returns the current MailSlurper version
*/
func GetVersion(writer http.ResponseWriter, request *http.Request) {
	var err error
	var result *server.Version
	log := context.Get(request, "log").(*logging.Logger)

	if result, err = server.GetServerVersionFromMaster(); err != nil {
		log.Errorf("Error getting version file from Github: %s", err.Error())
		GoHttpService.Error(writer, "There was an error reading the version file from GitHub")
		return
	}

	GoHttpService.WriteJson(writer, result, 200)
}
