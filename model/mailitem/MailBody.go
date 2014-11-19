// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailitem

import (
	"fmt"
	"strings"

	"github.com/adampresley/mailslurper/libmailslurper/model/attachment"
)

type MailBody struct {
	TextBody    string
	HTMLBody    string
	Attachments []*attachment.Attachment
}

func NewMailBody(textBody, htmlBody string, attachments []*attachment.Attachment) *MailBody {
	return &MailBody{
		TextBody:    textBody,
		HTMLBody:    htmlBody,
		Attachments: attachments,
	}
}

/*
Parses a mail's DATA section. This will attempt to figure out
what this mail contains. At the simples level it will contain
a text message. A more complex example would be a multipart message
with mixed text and HTML. It will also parse any attachments and
retrieve their contents into an attachments array.
*/
func (this *MailBody) Parse(contents string, boundary string) {
	/*
	 * Split the DATA content by CRLF CRLF. The first item will be the data
	 * headers. Everything past that is body/message.
	 */
	headerBodySplit := strings.Split(contents, "\r\n\r\n")
	if len(headerBodySplit) < 2 {
		panic("Expected DATA block to contain a header section and a body section")
	}

	contents = strings.Join(headerBodySplit[1:], "\r\n\r\n")
	this.Attachments = make([]*attachment.Attachment, 0)

	/*
	 * If there is no boundary then this is the simplest
	 * plain text type of mail you can get.
	 */
	if len(boundary) <= 0 {
		this.TextBody = contents
	} else {
		bodyParts := strings.Split(strings.TrimSpace(contents), fmt.Sprintf("--%s", strings.TrimSpace(boundary)))
		var index int

		/*
		 * First parse the headers for each of these attachments, then
		 * place each where they go.
		 */
		for index = 0; index < len(bodyParts); index++ {
			if len(strings.TrimSpace(bodyParts[index])) <= 0 || strings.TrimSpace(bodyParts[index]) == "--" {
				continue
			}

			header := &attachment.AttachmentHeader{}
			header.Parse(strings.TrimSpace(bodyParts[index]))

			switch {
			case strings.Contains(header.ContentType, "text/plain"):
				this.TextBody = header.Body

			case strings.Contains(header.ContentType, "text/html"):
				this.HTMLBody = header.Body

			case strings.Contains(header.ContentDisposition, "attachment"):
				newAttachment := &attachment.Attachment{
					Headers:  header,
					Contents: header.Body,
				}

				this.Attachments = append(this.Attachments, newAttachment)
			}
		}
	}
}
