package storage

import (
	"github.com/mailslurper/libmailslurper/model/attachment"
	"github.com/mailslurper/libmailslurper/model/mailitem"
	"github.com/mailslurper/libmailslurper/model/search"
)

/*
IStorage defines an interface for structures that need to connect to
storage engines. They store and retrieve data for MailSlurper
*/
type IStorage interface {
	Connect() error
	Disconnect()
	Create() error

	GetAttachment(mailID, attachmentID string) (attachment.Attachment, error)
	GetMailByID(id string) (mailitem.MailItem, error)
	GetMailCollection(offset, length int, mailSearch *search.MailSearch) ([]mailitem.MailItem, error)
	GetMailCount(mailSearch *search.MailSearch) (int, error)

	DeleteMailsAfterDate(startDate string) error
	StoreMail(mailItem *mailitem.MailItem) (string, error)
}
