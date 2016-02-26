// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package requests

import "github.com/mailslurper/libmailslurper/model/seed"

/*
DeleteMailRequest is used when requesting to delete mail
items.
*/
type DeleteMailRequest struct {
	PruneCode seed.PruneCode `json:"pruneCode"`
}
