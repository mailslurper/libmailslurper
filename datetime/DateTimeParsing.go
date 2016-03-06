// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package datetime

import (
	"log"
	"strings"
	"time"
)

var dateInputFormats = []string{
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -0700 MST",
	"Mon, 02 Jan 2006 15:04:05 -0700 (MST)",
	"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
	"02 Jan 2006 15:04:05 -0700",
}

/*
ParseDateTime takes a date/time string and attempts to parse it and return a newly formatted
date/time that looks like YYYY-MM-DD HH:MM:SS
*/
func ParseDateTime(dateString string) string {
	outputFormat := "2006-01-02 15:04:05"
	var parsedTime time.Time
	var err error

	dateString = strings.TrimSpace(dateString)
	result := ""

	for _, inputFormat := range dateInputFormats {
		if parsedTime, err = time.Parse(inputFormat, dateString); err == nil {
			result = parsedTime.Format(outputFormat)
			break
		}
	}

	if result == "" {
		log.Printf("libmailslurper: ERROR - Parsing date %s", dateString)
		result = dateString
	}

	return result
}
