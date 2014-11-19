// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package datetime

import (
	"log"
	"strings"
	"time"
)

/*
Takes a date/time string and attempts to parse it and return a newly formatted
date/time that looks like YYYY-MM-DD HH:MM:SS
*/
func ParseDateTime(dateString string) string {
	outputForm := "2006-01-02 15:04:05"
	firstForm := "Mon, 02 Jan 2006 15:04:05 -0700 MST"
	secondForm := "Mon, 02 Jan 2006 15:04:05 -0700 (MST)"
	thirdForm := "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"

	dateString = strings.TrimSpace(dateString)
	result := ""

	t, err := time.Parse(firstForm, dateString)
	if err != nil {
		t, err = time.Parse(secondForm, dateString)
		if err != nil {
			t, err = time.Parse(thirdForm, dateString)
			if err != nil {
				log.Printf("Error parsing date: %s\n", err)
				result = dateString
			} else {
				result = t.Format(outputForm)
			}
		} else {
			result = t.Format(outputForm)
		}
	} else {
		result = t.Format(outputForm)
	}

	return result
}
