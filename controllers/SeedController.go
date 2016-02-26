// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/mailslurper/libmailslurper/model/seed"
)

/*
GetPruneOptions returns a set of valid pruning options.

	GET: /v1/pruneoptions
*/
func GetPruneOptions(writer http.ResponseWriter, request *http.Request) {
	GoHttpService.WriteJson(writer, seed.PruneOptions, 200)
}
