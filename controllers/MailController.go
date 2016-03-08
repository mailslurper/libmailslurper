// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/adampresley/GoHttpService"
	"github.com/adampresley/logging"
	"github.com/mailslurper/libmailslurper/model/mailitem"
	"github.com/mailslurper/libmailslurper/model/requests"
	"github.com/mailslurper/libmailslurper/model/response"
	"github.com/mailslurper/libmailslurper/model/search"
	"github.com/mailslurper/libmailslurper/storage"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

/*
DeleteMail is a request to delete mail items. This expects a body containing
a DeleteMailRequest object.

	DELETE: /mail
*/
func DeleteMail(writer http.ResponseWriter, request *http.Request) {
	var err error
	deleteRequest := &requests.DeleteMailRequest{}

	log := context.Get(request, "log").(*logging.Logger)
	database := context.Get(request, "database").(storage.IStorage)

	if err = GoHttpService.ParseJsonBody(request, deleteRequest); err != nil {
		log.Errorf("Invalid delete mail request - %s", err.Error())
		GoHttpService.BadRequest(writer, "Invalid delete mail request")
		return
	}

	if !deleteRequest.PruneCode.IsValid() {
		log.Errorf("Attempt to use invalid prune code - %s", deleteRequest.PruneCode)
		GoHttpService.BadRequest(writer, "Invalid prune type")
		return
	}

	startDate := deleteRequest.PruneCode.ConvertToDate()

	if err = database.DeleteMailsAfterDate(startDate); err != nil {
		log.Errorf("Problem deleting mails - %s", err.Error())
		GoHttpService.Error(writer, "There was a problem deleting mails")
		return
	}

	log.Infof("Deleting mails, code %s - Start - %s", deleteRequest.PruneCode.String(), startDate)
	GoHttpService.Success(writer, "OK")
}

/*
GetMail returns a single mail item by ID.

	GET: /mail/{id}
*/
func GetMail(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	log := context.Get(request, "log").(*logging.Logger)
	database := context.Get(request, "database").(storage.IStorage)

	var mailID string
	var mailItem mailitem.MailItem
	var err error
	var ok bool

	/*
	 * Validate incoming arguments
	 */
	if mailID, ok = vars["mailID"]; !ok {
		log.Error("Invalid mail ID passed to GetMail")
		GoHttpService.BadRequest(writer, "A valid mail ID is required")
		return
	}

	/*
	 * Retrieve the mail item
	 */
	if mailItem, err = database.GetMailByID(mailID); err != nil {
		log.Errorf("Problem getting mail item in GetMail - %s", err.Error())
		GoHttpService.Error(writer, "Problem getting mail item")
		return
	}

	log.Infof("Mail item %s retrieved", mailID)

	result := &response.MailItemResponse{
		MailItem: mailItem,
	}

	GoHttpService.WriteJson(writer, result, 200)
}

/*
GetMailCollection returns a collection of mail items. This is constrianed
by a page number. A page of data contains 50 items.

	GET: /mails?pageNumber={pageNumber}
*/
func GetMailCollection(writer http.ResponseWriter, request *http.Request) {
	var err error
	var pageNumberString string
	var pageNumber int
	var mailCollection []mailitem.MailItem
	var totalRecordCount int

	log := context.Get(request, "log").(*logging.Logger)
	database := context.Get(request, "database").(storage.IStorage)

	/*
	 * Validate incoming arguments. A page is currently 50 items, hard coded
	 */
	pageNumberString = request.URL.Query().Get("pageNumber")
	if pageNumberString == "" {
		pageNumber = 1
	} else {
		if pageNumber, err = strconv.Atoi(pageNumberString); err != nil {
			log.Error("Invalid page number passed to GetMailCollection")
			GoHttpService.BadRequest(writer, "A valid page number is required")
			return
		}
	}

	length := 50
	offset := (pageNumber - 1) * length

	/*
	 * Retrieve mail items
	 */
	mailSearch := &search.MailSearch{
		Message: request.URL.Query().Get("message"),
		Start:   request.URL.Query().Get("start"),
		End:     request.URL.Query().Get("end"),
		From:    request.URL.Query().Get("from"),
		To:      request.URL.Query().Get("to"),

		OrderByField:     request.URL.Query().Get("orderby"),
		OrderByDirection: request.URL.Query().Get("dir"),
	}

	if mailCollection, err = database.GetMailCollection(offset, length, mailSearch); err != nil {
		log.Errorf("Problem getting mail collection - %s", err.Error())
		GoHttpService.Error(writer, "Problem getting mail collection")
		return
	}

	if totalRecordCount, err = database.GetMailCount(mailSearch); err != nil {
		log.Errorf("Problem getting record count in GetMailCollection - %s", err.Error())
		GoHttpService.Error(writer, "Error getting record count")
		return
	}

	totalPages := int(math.Ceil(float64(totalRecordCount / length)))
	if totalPages*length < totalRecordCount {
		totalPages++
	}

	log.Infof("Mail collection page %d retrieved", pageNumber)

	result := &response.MailCollectionResponse{
		MailItems:    mailCollection,
		TotalPages:   totalPages,
		TotalRecords: totalRecordCount,
	}

	GoHttpService.WriteJson(writer, result, 200)
}

/*
GetMailCount returns the number of mail items in storage.

	GET: /mailcount
*/
func GetMailCount(writer http.ResponseWriter, request *http.Request) {
	var err error
	var mailItemCount int

	log := context.Get(request, "log").(*logging.Logger)
	database := context.Get(request, "database").(storage.IStorage)

	/*
	 * Get the count
	 */
	if mailItemCount, err = database.GetMailCount(&search.MailSearch{}); err != nil {
		log.Errorf("Problem getting mail item count in GetMailCount - %s", err.Error())
		GoHttpService.Error(writer, "Problem getting mail count")
		return
	}

	log.Infof("Mail item count - %d", mailItemCount)

	result := response.MailCountResponse{
		MailCount: mailItemCount,
	}

	GoHttpService.WriteJson(writer, result, 200)
}
