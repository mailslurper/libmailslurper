// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/adampresley/GoHttpService"
	"github.com/adampresley/logging"
	"github.com/mailslurper/libmailslurper/model/attachment"
	"github.com/mailslurper/libmailslurper/storage"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

/*
DownloadAttachment retrieves binary database from storage and streams
it back to the caller
*/
func DownloadAttachment(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	log := context.Get(request, "log").(*logging.Logger)
	database := context.Get(request, "database").(storage.IStorage)

	var err error
	var attachmentID string
	var mailID string
	var ok bool

	var attachment attachment.Attachment
	var data []byte

	/*
	 * Validate incoming arguments
	 */
	if mailID, ok = vars["mailID"]; !ok {
		log.Error("No valid mail ID passed to DownloadAttachment")
		GoHttpService.BadRequest(writer, "A valid mail ID is required")
		return
	}

	if attachmentID, ok = vars["attachmentId"]; !ok {
		log.Error("No valid attachment ID passed to DownloadAttachment")
		GoHttpService.BadRequest(writer, "A valid attachment ID is required")
		return
	}

	/*
	 * Retrieve the attachment
	 */
	if attachment, err = database.GetAttachment(mailID, attachmentID); err != nil {
		log.Errorf("Problem getting attachment %s - %s", attachmentID, err.Error())
		GoHttpService.Error(writer, "Error getting attachment "+attachmentID)
		return
	}

	/*
	 * Decode the base64 data and stream it back
	 */
	if attachment.IsContentBase64() {
		data, err = base64.StdEncoding.DecodeString(attachment.Contents)
		if err != nil {
			log.Errorf("Problem decoding attachment %s - %s", attachmentID, err.Error())
			GoHttpService.Error(writer, "Cannot decode attachment")
			return
		}
	} else {
		data = []byte(attachment.Contents)
	}

	log.Infof("Attachment %s retrieved", attachmentID)

	reader := bytes.NewReader(data)
	http.ServeContent(writer, request, attachment.Headers.FileName, time.Now(), reader)
}
