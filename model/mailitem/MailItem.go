// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailitem

import (
	"github.com/adampresley/mailslurper/libmailslurper/model/attachment"
)

/*
MailItem is a struct describing a parsed mail item. This is
populated after an incoming client connection has finished
sending mail data to this server.
*/
type MailItem struct {
	Id          string                   `json:"id"`
	DateSent    string                   `json:"dateSent"`
	FromAddress string                   `json:"fromAddress"`
	ToAddresses []string                 `json:"toAddresses"`
	Subject     string                   `json:"subject"`
	XMailer     string                   `json:"xmailer"`
	Body        string                   `json:"body"`
	ContentType string                   `json:"contentType"`
	Boundary    string                   `json:"boundary"`
	Attachments []*attachment.Attachment `json:"attachments"`
}

func NewMailItem(id, dateSent, fromAddress string, toAddresses []string, subject, xMailer, body, contentType, boundary string, attachments []*attachment.Attachment) *MailItem {
	return &MailItem{
		Id:          id,
		DateSent:    dateSent,
		FromAddress: fromAddress,
		ToAddresses: toAddresses,
		Subject:     subject,
		XMailer:     xMailer,
		Body:        body,
		ContentType: contentType,
		Boundary:    boundary,
		Attachments: attachments,
	}
}
