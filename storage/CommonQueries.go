package storage

func getMailAndAttachmentsQuery(whereClause string) string {
	sqlQuery := `
		SELECT
			  mailitem.dateSent
			, mailitem.fromAddress
			, mailitem.toAddressList
			, mailitem.subject
			, mailitem.xmailer
			, mailitem.body
			, mailitem.contentType
			, mailitem.boundary
			, attachment.id AS attachmentID
			, attachment.fileName
			, attachment.contentType

		FROM mailitem
			LEFT JOIN attachment ON attachment.mailItemID=mailitem.id

		WHERE 1=1 `

	sqlQuery = sqlQuery + whereClause
	sqlQuery = sqlQuery + ` ORDER BY mailitem.dateSent DESC `

	return sqlQuery
}

func getMailCountQuery() string {
	sqlQuery := `
		SELECT COUNT(id) AS mailItemCount FROM mailitem
	`

	return sqlQuery
}

func getDeleteMailQuery(startDate string) string {
	where := ""

	if len(startDate) > 0 {
		where = where + " AND dateSent <= ? "
	}

	sqlQuery := "DELETE FROM mailitem WHERE 1=1" + where
	return sqlQuery
}

func getInsertMailQuery() string {
	sqlQuery := `
		INSERT INTO mailitem (
			  id
			, dateSent
			, fromAddress
			, toAddressList
			, subject
			, xmailer
			, body
			, contentType
			, boundary
		) VALUES (
			  ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
		)
	`

	return sqlQuery
}

func getInsertAttachmentQuery() string {
	sqlQuery := `
		INSERT INTO attachment (
			  id
			, mailItemId
			, fileName
			, contentType
			, content
		) VALUES (
			  ?
			, ?
			, ?
			, ?
			, ?
		)
	`

	return sqlQuery
}
