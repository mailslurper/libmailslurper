// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package receiver

import (
	"log"

	"github.com/mailslurper/libmailslurper/model/mailitem"
	"github.com/mailslurper/libmailslurper/storage"
)

type DatabaseReceiver struct{}

func (this DatabaseReceiver) Receive(mailItem *mailitem.MailItem) error {
	newId, err := storage.StoreMail(mailItem)
	if err != nil {
		log.Println("libmailslurper: ERROR - There was an error while storing your mail item:", err)
		return err
	}

	log.Println("libmailslurper: INFO - Mail item", newId, "written")
	return nil
}
