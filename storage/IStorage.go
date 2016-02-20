package storage

import (
	"github.com/mailslurper/libmailslurper/model/attachment"
	"github.com/mailslurper/libmailslurper/model/mailitem"
)

/*
IStorage defines an interface for structures that need to connect to
storage engines. They store and retrieve data for MailSlurper
*/
type IStorage interface {
	Connect() error
	Disconnect()

	GetAttachment(mailID, attachmentID string) (attachment.Attachment, error)
	GetMailByID(id string) (mailitem.MailItem, error)
	GetMailCollection(offset, length int) ([]mailitem.MailItem, error)
	GetMailCount() (int, error)

	DeleteMailsAfterDate(startDate string) error
	StoreMail(mailItem *mailitem.MailItem) (string, error)
}
