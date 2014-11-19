// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package attachment

import (
	"log"
	"strings"

	"github.com/adampresley/mailslurper/libmailslurper/model/header"
)

type AttachmentHeader struct {
	ContentType             string `json:"contentType"`
	MIMEVersion             string `json:"mimeVersion"`
	ContentTransferEncoding string `json:"contentTransferEncoding"`
	ContentDisposition      string `json:"contentDisposition"`
	FileName                string `json:"fileName"`
	Body                    string `json:"body"`
}

func NewAttachmentHeader(contentType, mimeVersion, contentTransferEncoding, contentDisposition, fileName, body string) *AttachmentHeader {
	return &AttachmentHeader{
		ContentType:             contentType,
		MIMEVersion:             mimeVersion,
		ContentTransferEncoding: contentTransferEncoding,
		ContentDisposition:      contentDisposition,
		FileName:                fileName,
		Body:                    body,
	}
}

/*
Parses a set of attachment headers. Splits lines up and figures out what
header data goes into what structure key. Most headers follow this format:

Header-Name: Some value here\r\n
*/
func (this *AttachmentHeader) Parse(contents string) {
	var key string

	headerBodySplit := strings.Split(contents, "\r\n\r\n")
	if len(headerBodySplit) < 2 {
		panic("Expected attachment to contain a header section and a body section")
	}

	contents = headerBodySplit[0]

	this.Body = strings.Join(headerBodySplit[1:], "\r\n\r\n")
	this.FileName = ""
	this.ContentType = ""
	this.ContentDisposition = ""
	this.ContentTransferEncoding = ""
	this.MIMEVersion = ""

	/*
	 * Unfold and split the header into lines. Loop over each line
	 * and figure out what headers are present. Store them.
	 * Sadly some headers require special processing.
	 */
	contents = header.UnfoldHeaders(contents)
	splitHeader := strings.Split(contents, "\r\n")
	numLines := len(splitHeader)

	for index := 0; index < numLines; index++ {
		splitItem := strings.Split(splitHeader[index], ":")
		key = splitItem[0]

		switch strings.ToLower(key) {
		case "content-disposition":
			contentDisposition := strings.TrimSpace(strings.Join(splitItem[1:], ""))
			log.Println("Attachment Content-Disposition: ", contentDisposition)

			contentDispositionSplit := strings.Split(contentDisposition, ";")
			contentDispositionRightSide := strings.TrimSpace(strings.Join(contentDispositionSplit[1:], ";"))

			if len(contentDispositionSplit) < 2 || (len(contentDispositionSplit) > 1 && len(strings.TrimSpace(contentDispositionRightSide)) <= 0) {
				this.ContentDisposition = contentDisposition
			} else {
				this.ContentDisposition = strings.TrimSpace(contentDispositionSplit[0])

				/*
				 * See if we have an attachment and filename
				 */
				if strings.Contains(strings.ToLower(this.ContentDisposition), "attachment") && len(strings.TrimSpace(contentDispositionRightSide)) > 0 {
					filenameSplit := strings.Split(contentDispositionRightSide, "=")
					this.FileName = strings.Replace(strings.Join(filenameSplit[1:], "="), "\"", "", -1)
				}
			}

		case "content-transfer-encoding":
			this.ContentTransferEncoding = strings.TrimSpace(strings.Join(splitItem[1:], ""))
			log.Println("Attachment Content-Transfer-Encoding: ", this.ContentTransferEncoding)

		case "content-type":
			contentType := strings.TrimSpace(strings.Join(splitItem[1:], ""))
			log.Println("Attachment Content-Type: ", contentType)

			contentTypeSplit := strings.Split(contentType, ";")

			if len(contentTypeSplit) < 2 {
				this.ContentType = contentType
			} else {
				this.ContentType = strings.TrimSpace(contentTypeSplit[0])
				contentTypeRightSide := strings.TrimSpace(strings.Join(contentTypeSplit[1:], ";"))

				/*
				 * See if there is a "name" portion to this
				 */
				if strings.Contains(strings.ToLower(contentTypeRightSide), "name") || strings.Contains(strings.ToLower(contentTypeRightSide), "filename") {
					filenameSplit := strings.Split(contentTypeRightSide, "=")
					this.FileName = strings.Replace(strings.Join(filenameSplit[1:], "="), "\"", "", -1)
				}
			}

		case "mime-version":
			this.MIMEVersion = strings.TrimSpace(strings.Join(splitItem[1:], ""))
			log.Println("Attachment MIME-Version: ", this.MIMEVersion)
		}
	}
}
