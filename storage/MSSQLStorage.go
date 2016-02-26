package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/adampresley/sanitizer"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/mailslurper/libmailslurper/model/attachment"
	"github.com/mailslurper/libmailslurper/model/mailitem"
	"github.com/mailslurper/libmailslurper/model/search"
)

/*
MSSQLStorage implements the IStorage interface
*/
type MSSQLStorage struct {
	connectionInformation *ConnectionInformation
	db                    *sql.DB
	xssService            sanitizer.XSSServiceProvider
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
			  attachment.fileName
			, attachment.contentType
			, attachment.content
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
func (storage *MSSQLStorage) GetMailByID(mailItemID string) (mailitem.MailItem, error) {
	result := mailitem.MailItem{}
	attachments := make([]*attachment.Attachment, 0)

	var err error
	var rows *sql.Rows

	var dateSent string
	var fromAddress string
	var toAddressList string
	var subject string
	var xmailer string
	var body string
	var boundary string
	var attachmentID sql.NullString
	var fileName string
	var contentType string

	sqlQuery := getMailAndAttachmentsQuery(" AND mailitem.id=? ")

	if rows, err = storage.db.Query(sqlQuery, mailItemID); err != nil {
		return result, fmt.Errorf("Error running query to get mail item: %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&dateSent, &fromAddress, &toAddressList, &subject, &xmailer, &body, &contentType, &boundary, &attachmentID, &fileName, &contentType)
		if err != nil {
			return result, fmt.Errorf("Error scanning mail item record: %s", err.Error())
		}

		/*
		 * Only capture the mail item once. Every subsequent record is an attachment
		 */
		if result.Id == "" {
			result = mailitem.MailItem{
				Id:          mailItemID,
				DateSent:    dateSent,
				FromAddress: fromAddress,
				ToAddresses: strings.Split(toAddressList, "; "),
				Subject:     storage.xssService.SanitizeString(subject),
				XMailer:     storage.xssService.SanitizeString(xmailer),
				Body:        storage.xssService.SanitizeString(body),
				ContentType: contentType,
				Boundary:    boundary,
			}
		}

		if attachmentID.Valid {
			newAttachment := &attachment.Attachment{
				Id:     attachmentID.String,
				MailId: mailItemID,
				Headers: &attachment.AttachmentHeader{
					FileName:    storage.xssService.SanitizeString(fileName),
					ContentType: contentType,
				},
			}

			attachments = append(attachments, newAttachment)
		}
	}

	result.Attachments = attachments
	return result, nil
}

/*
GetMailCollection retrieves a slice of mail items starting at offset and getting length number
of records. This query is MSSQL 2005 and higher compatible.
*/
func (storage *MSSQLStorage) GetMailCollection(offset, length int, mailSearch *search.MailSearch) ([]mailitem.MailItem, error) {
	result := make([]mailitem.MailItem, 0)
	attachments := make([]*attachment.Attachment, 0)

	var err error
	var rows *sql.Rows

	var currentMailItemID string
	var currentMailItem mailitem.MailItem
	var parameters []interface{}

	var mailItemID string
	var dateSent string
	var fromAddress string
	var toAddressList string
	var subject string
	var xmailer string
	var body string
	var mailContentType string
	var boundary string
	var attachmentID sql.NullString
	var fileName sql.NullString
	var attachmentContentType sql.NullString

	/*
	 * This query is MSSQL 2005 and higher compatible
	 */
	sqlQuery := `
		WITH pagedMailItems AS (
			SELECT
				  mailitem.id
				, mailitem.dateSent
				, mailitem.fromAddress
				, mailitem.toAddressList
				, mailitem.subject
				, mailitem.xmailer
				, mailitem.body
				, mailitem.contentType
				, mailitem.boundary
				, ROW_NUMBER() OVER (ORDER BY mailitem.dateSent DESC) AS rowNumber
			FROM mailitem
			WHERE 1=1
	`

	sqlQuery, parameters = addSearchCriteria(sqlQuery, parameters, mailSearch)

	sqlQuery = sqlQuery + `
		)
		SELECT
			  pagedMailItems.id AS mailItemID
			, pagedMailItems.dateSent
			, pagedMailItems.fromAddress
			, pagedMailItems.toAddressList
			, pagedMailItems.subject
			, pagedMailItems.xmailer
			, pagedMailItems.body
			, pagedMailItems.contentType AS mailContentType
			, pagedMailItems.boundary
			, attachment.id AS attachmentID
			, attachment.fileName
			, attachment.contentType AS attachmentContentType

		FROM pagedMailItems
			LEFT JOIN attachment ON attachment.mailItemID=pagedMailItems.id

		WHERE
			pagedMailItems.rowNumber BETWEEN ? AND ?

		ORDER BY pagedMailItems.dateSent DESC
	`

	parameters = append(parameters, offset)
	parameters = append(parameters, offset+length)

	if rows, err = storage.db.Query(sqlQuery, parameters...); err != nil {
		return result, fmt.Errorf("Error running query to get mail collection: %s", err.Error())
	}

	defer rows.Close()

	currentMailItemID = ""

	for rows.Next() {
		err = rows.Scan(&mailItemID, &dateSent, &fromAddress, &toAddressList, &subject, &xmailer, &body, &mailContentType, &boundary, &attachmentID, &fileName, &attachmentContentType)
		if err != nil {
			return result, fmt.Errorf("Error scanning mail item record: %s", err.Error())
		}

		if currentMailItemID != mailItemID {
			/*
			 * If we have a mail item we are working with place the attachments with it.
			 * Then reset everything in prep for the next mail item and batch of attachments
			 */
			if currentMailItemID != "" {
				currentMailItem.Attachments = attachments
				result = append(result, currentMailItem)
			}

			currentMailItem = mailitem.MailItem{
				Id:          mailItemID,
				DateSent:    dateSent,
				FromAddress: fromAddress,
				ToAddresses: strings.Split(toAddressList, "; "),
				Subject:     storage.xssService.SanitizeString(subject),
				XMailer:     storage.xssService.SanitizeString(xmailer),
				Body:        storage.xssService.SanitizeString(body),
				ContentType: mailContentType,
				Boundary:    boundary,
			}

			currentMailItemID = mailItemID
			attachments = make([]*attachment.Attachment, 0)
		}

		if attachmentID.Valid {
			newAttachment := &attachment.Attachment{
				Id:     attachmentID.String,
				MailId: mailItemID,
				Headers: &attachment.AttachmentHeader{
					FileName:    storage.xssService.SanitizeString(fileName.String),
					ContentType: attachmentContentType.String,
				},
			}

			attachments = append(attachments, newAttachment)
		}
	}

	/*
	 * Attach our straggler
	 */
	if currentMailItemID != "" {
		currentMailItem.Attachments = attachments
		result = append(result, currentMailItem)
	}

	return result, nil
}

/*
GetMailCount returns the number of total records in the mail items table
*/
func (storage *MSSQLStorage) GetMailCount(mailSearch *search.MailSearch) (int, error) {
	var mailItemCount int
	var err error

	sqlQuery, parameters := getMailCountQuery(mailSearch)
	if err = storage.db.QueryRow(sqlQuery, parameters...).Scan(&mailItemCount); err != nil {
		return 0, fmt.Errorf("Error running query to get mail item count: %s", err.Error())
	}

	return mailItemCount, nil
}

/*
DeleteMailsAfterDate deletes all mails after a specified date
*/
func (storage *MSSQLStorage) DeleteMailsAfterDate(startDate string) error {
	sqlQuery := getDeleteMailQuery(startDate)
	parameters := []interface{}{}
	var err error

	if len(startDate) > 0 {
		parameters = append(parameters, startDate)
	}

	_, err = storage.db.Exec(sqlQuery, parameters...)
	return err
}

/*
StoreMail writes a mail item and its attachments to the storage device. This returns the new mail ID
*/
func (storage *MSSQLStorage) StoreMail(mailItem *mailitem.MailItem) (string, error) {
	var err error
	var transaction *sql.Tx
	var statement *sql.Stmt

	/*
	 * Create a transaction and insert the new mail item
	 */
	if transaction, err = storage.db.Begin(); err != nil {
		return "", fmt.Errorf("Error starting transaction in StoreMail: %s", err.Error())
	}

	/*
	 * Insert the mail item
	 */
	if statement, err = transaction.Prepare(getInsertMailQuery()); err != nil {
		return "", fmt.Errorf("Error preparing insert statement for mail item in StoreMail: %s", err.Error())
	}

	_, err = statement.Exec(
		mailItem.Id,
		mailItem.DateSent,
		mailItem.FromAddress,
		strings.Join(mailItem.ToAddresses, "; "),
		mailItem.Subject,
		mailItem.XMailer,
		mailItem.Body,
		mailItem.ContentType,
		mailItem.Boundary,
	)

	if err != nil {
		transaction.Rollback()
		return "", fmt.Errorf("Error executing insert for mail item in StoreMail: %s", err.Error())
	}

	statement.Close()

	/*
	 * Insert attachments
	 */
	if err = storeAttachments(mailItem.Id, transaction, mailItem.Attachments); err != nil {
		transaction.Rollback()
		return "", fmt.Errorf("Unable to insert attachments in StoreMail: %s", err.Error())
	}

	transaction.Commit()
	log.Printf("New mail item written to database.\n\n")

	return mailItem.Id, nil
}

func storeAttachments(mailItemID string, transaction *sql.Tx, attachments []*attachment.Attachment) error {
	var err error
	var attachmentID string

	for _, currentAttachment := range attachments {
		if attachmentID, err = mailitem.GenerateId(); err != nil {
			return fmt.Errorf("Error generating ID for attachment: %s", err.Error())
		}

		statement, err := transaction.Prepare(getInsertAttachmentQuery())
		if err != nil {
			return fmt.Errorf("Error preparing insert attachment statement: %s", err.Error())
		}

		_, err = statement.Exec(
			attachmentID,
			mailItemID,
			currentAttachment.Headers.FileName,
			currentAttachment.Headers.ContentType,
			currentAttachment.Contents,
		)

		if err != nil {
			return fmt.Errorf("Error executing insert attachment in StoreMail: %s", err.Error())
		}

		statement.Close()
		currentAttachment.Id = attachmentID
	}

	return nil
}
