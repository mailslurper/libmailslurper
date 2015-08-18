// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package attachment

import (
	"math"
	"regexp"
	"strings"
)

type Attachment struct {
	Id       string            `json:"id"`
	MailId   string            `json:"mailId"`
	Headers  *AttachmentHeader `json:"headers"`
	Contents string            `json:"contents"`
}

func NewAttachment(headers *AttachmentHeader, contents string) *Attachment {
	return &Attachment{
		Headers:  headers,
		Contents: contents,
	}
}

/*
IsContentBase64 returns true/false if the content of this attachment
resembles a base64 encoded string.
*/
func (attachment *Attachment) IsContentBase64() bool {
	spaceKiller := func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}

		return r
	}

	trimmedContents := strings.Map(spaceKiller, attachment.Contents)

	if math.Mod(float64(len(trimmedContents)), 4.0) == 0 {
		matchResult, err := regexp.Match("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$", []byte(trimmedContents))
		if err == nil {
			if matchResult {
				return true
			}
		}
	}

	return false
}
