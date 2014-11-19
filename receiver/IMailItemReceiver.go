// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package receiver

import(
	"github.com/adampresley/mailslurper/libmailslurper/model/mailitem"
)

type IMailItemReceiver interface{
	Receive(mailItem *mailitem.MailItem) error
}
