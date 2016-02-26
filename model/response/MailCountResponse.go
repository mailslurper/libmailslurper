// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package response

/*
MailCountResponse is used to report the number of mail
items in storage.
*/
type MailCountResponse struct {
	MailCount int `json:"mailCount"`
}
