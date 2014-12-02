// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailitem

import (
	"github.com/nu7hatch/gouuid"
)

/*
Generate a UUID ID for database records.
*/
func GenerateId() (string, error) {
	id, err := uuid.NewV4()
	return id.String(), err
}

