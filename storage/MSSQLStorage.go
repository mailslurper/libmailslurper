package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/adampresley/golangdb"
	"github.com/adampresley/sanitizer"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/mailslurper/libmailslurper/model/attachment"
	"github.com/mailslurper/libmailslurper/model/mailitem"
)

/*
MSSQLStorage implements the IStorage interface
*/
type MSSQLStorage struct {
	connectionInformation *ConnectionInformation
	db                    *sql.DB
	xssService            *sanitizer.XSSServiceProvider
}

/*
NewMSSQLStorage creates a new storage object that interfaces to MSSQL
*/
func NewMSSQLStorage(connectionInformation *ConnectionInformation) *MSSQLStorage {
	return &MSSQLStorage{
		connectionInformation: connectionInformation,
		xssService:            sanitizer.NewXSSService(),
	}
}

/*
Connect to the database
*/
func (storage *MSSQLStorage) Connect() error {
	connectionString := fmt.Sprintf("Server=%s;Port=%d;User Id=%s;Password=%s;Database=%s",
		storage.connectionInformation.Address,
		storage.connectionInformation.Port,
		storage.connectionInformation.UserName,
		storage.connectionInformation.Password,
		storage.connectionInformation.Database,
	)

	db, err := sql.Open("mssql", connectionString)

	storage.db = db
	return err
}

/*
Disconnect does exactly what you think it does
*/
func (storage *MSSQLStorage) Disconnect() {
	storage.db.Close()
}

/*
GetAttachment retrieves an attachment for a given mail item
*/
func (storage *MSSQLStorage) GetAttachment(mailID, attachmentID string) (attachment.Attachment, error) {
	result := attachment.Attachment{}
	var err error
	var rows *sql.Rows

	var fileName string
	var contentType string
	var content string

	getAttachmentSQL := `
		SELECT
			  fileName
			, contentType
			, content
		FROM attachment
		WHERE
			id=?
			AND mailItemId=?
	`

	if rows, err = storage.db.Query(getAttachmentSQL, attachmentID, mailID); err != nil {
		return result, fmt.Errorf("Error running query to get attachment")
	}

	defer rows.Close()
	rows.Next()
	rows.Scan(&fileName, &contentType, &content)

	result.Headers = &attachment.AttachmentHeader{
		FileName:    fileName,
		ContentType: contentType,
	}

	result.Contents = content
	return result, nil
}

/*
GetMailByID retrieves a single mail item and attachment by ID
*/
func (storage *MSSQLStorage) GetMailByID(id string) (mailitem.MailItem, error) {
	result := mailitem.MailItem{}

	sql := getMailQuery(" AND mailitem.id=? ")
	rows, err := golangdb.Db["lib"].Query(sql, id)

	if err != nil {
		return result, fmt.Errorf("Error running query to get mail items: %s", err)
	}

	rows.Next()

	var mailItemId string
	var dateSent string
	var fromAddress string
	var toAddressList string
	var subject string
	var xmailer string
	var body string
	var contentType string
	var boundary string

	rows.Scan(&mailItemId, &dateSent, &fromAddress, &toAddressList, &subject, &xmailer, &body, &contentType, &boundary)

	result = mailitem.MailItem{
		Id:          mailItemId,
		DateSent:    dateSent,
		FromAddress: fromAddress,
		ToAddresses: strings.Split(toAddressList, "; "),
		Subject:     xssService.SanitizeString(subject),
		XMailer:     xssService.SanitizeString(xmailer),
		Body:        xssService.SanitizeString(body),
		ContentType: contentType,
		Boundary:    boundary,
	}

	/*
	 * Get attachments
	 */
	sql = `
		SELECT
			  attachment.id AS attachmentId
			, attachment.fileName
			, attachment.contentType

		FROM attachment
		WHERE
			attachment.mailItemId=?`

	attachmentRows, err := golangdb.Db["lib"].Query(sql, mailItemId)
	if err != nil {
		return result, err
	}

	attachments := make([]*attachment.Attachment, 0)

	for attachmentRows.Next() {
		var attachmentId string
		var fileName string
		var contentType string

		attachmentRows.Scan(&attachmentId, &fileName, &contentType)

		newAttachment := &attachment.Attachment{
			Id:     attachmentId,
			MailId: mailItemId,
			Headers: &attachment.AttachmentHeader{
				FileName:    xssService.SanitizeString(fileName),
				ContentType: contentType,
			},
		}

		attachments = append(attachments, newAttachment)
	}

	attachmentRows.Close()

	result.Attachments = attachments

	rows.Close()
	return result, nil
}
