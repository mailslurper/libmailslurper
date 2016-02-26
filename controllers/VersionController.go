// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"fmt"
	"net/http"

	"github.com/adampresley/GoHttpService"
)

/*
GetVersion returns the current MailSlurper version
*/
func GetVersion(writer http.ResponseWriter, request *http.Request) {
	GoHttpService.Success(writer, fmt.Sprintf("MailSlurperService server version 1.8"))
}
