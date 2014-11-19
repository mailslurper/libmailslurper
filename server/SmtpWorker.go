// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package server

import (
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/adampresley/mailslurper/libmailslurper/model/attachment"
	"github.com/adampresley/mailslurper/libmailslurper/model/mailitem"
	"github.com/adampresley/mailslurper/libmailslurper/smtpconstants"
	"github.com/adampresley/mailslurper/libmailslurper/smtpio"
)

type SmtpWorker struct{
	Connection *net.TCPConn
	Mail       mailitem.MailItem
	Reader     smtpio.SmtpReader
	Receiver   chan mailitem.MailItem
	State      smtpconstants.SmtpWorkerState
	WorkerId   int
	Writer     smtpio.SmtpWriter
}

/*
This function takes a command and the raw data read from the socket
connection and executes the correct handler function to process
the data and potentially respond to the client to continue SMTP negotiations.
*/
func (this *SmtpWorker) ExecuteCommand(command smtpconstants.SmtpCommand, streamInput string) error {
	var err error
	var response string

	var headers mailitem.MailHeader
	var body mailitem.MailBody

	streamInput = strings.TrimSpace(streamInput)

	switch command {
	case smtpconstants.HELO:
		err = this.Process_HELO(streamInput)

	case smtpconstants.MAIL:
		response, err := this.Process_MAIL(streamInput)
		if err != nil {
			log.Println("ERROR -", err)
		} else {
			this.Mail.FromAddress = response
			log.Println("Mail from", response)
		}

	case smtpconstants.RCPT:
		response, err = this.Process_RCPT(streamInput)
		if err != nil {
			log.Println("ERROR -", err)
		} else {
			this.Mail.ToAddresses = append(this.Mail.ToAddresses, response)
		}

	case smtpconstants.DATA:
		headers, body, err = this.Process_DATA(streamInput)
		if err != nil {
			log.Println("ERROR -", err)
		} else {
			if len(strings.TrimSpace(body.HTMLBody)) <= 0 {
				this.Mail.Body = body.TextBody
			} else {
				this.Mail.Body = body.HTMLBody
			}

			this.Mail.Subject = headers.Subject
			this.Mail.DateSent = headers.Date
			this.Mail.XMailer = headers.XMailer
			this.Mail.ContentType = headers.ContentType
			this.Mail.Boundary = headers.Boundary
			this.Mail.Attachments = body.Attachments
		}

	default:
		err = nil
	}

	return err
}

/*
Initializes the mail item structure that will eventually
be written to the receiver channel.
*/
func (this *SmtpWorker) InitializeMailItem() {
	this.Mail.ToAddresses = make([]string, 0)
	this.Mail.Attachments = make([]*attachment.Attachment, 0)
}

/*
Function to process the DATA command (constant DATA). When a client sends the DATA
command there are three parts to the transmission content. Before this data
can be processed this function will tell the client how to terminate the DATA block.
We are asking clients to terminate with "\r\n.\r\n".

The first part is a set of header lines. Each header line is a header key (name), followed
by a colon, followed by the value for that header key. For example a header key might
be "Subject" with a value of "Testing Mail!".

After the header section there should be two sets of carriage return/line feed characters.
This signals the end of the header block and the start of the message body.

Finally when the client sends the "\r\n.\r\n" the DATA transmission portion is complete.
This function will return the following items.

	1. Headers (MailHeader)
	2. Body breakdown (MailBody)
	3. error structure
*/
func (this *SmtpWorker) Process_DATA(streamInput string) (mailitem.MailHeader, mailitem.MailBody, error) {
	header := mailitem.MailHeader{}
	body := mailitem.MailBody{}

	commandCheck := strings.Index(strings.ToLower(streamInput), "data")
	if commandCheck < 0 {
		return header, body, errors.New("Invalid command for DATA")
	}

	this.Writer.SendDataResponse()
	entireMailContents := this.Reader.ReadDataBlock()

	/*
	 * Parse the header content
	 */
	header.Parse(entireMailContents)

	/*
	 * Parse the body
	 */
	body.Parse(entireMailContents, header.Boundary)

	this.Writer.SendOkResponse()
	return header, body, nil
}

/*
Function to process the HELO and EHLO SMTP commands. This command
responds to clients with a 250 greeting code and returns success
or false and an error message (if any).
*/
func (this *SmtpWorker) Process_HELO(streamInput string) error {
	lowercaseStreamInput := strings.ToLower(streamInput)

	commandCheck := (strings.Index(lowercaseStreamInput, "helo") + 1) + (strings.Index(lowercaseStreamInput, "ehlo") + 1)
	if commandCheck < 0 {
		return errors.New("Invalid HELO command")
	}

	split := strings.Split(streamInput, " ")
	if len(split) < 2 {
		return errors.New("HELO command format is invalid")
	}

	return this.Writer.SendHELOResponse()
}

/*
Function to process the MAIL FROM command (constant MAIL). This command
will respond to clients with 250 Ok response and returns an error
that may have occurred as well as the parsed FROM.
*/
func (this *SmtpWorker) Process_MAIL(streamInput string) (string, error) {
	commandCheck := strings.Index(strings.ToLower(streamInput), "mail from")
	if commandCheck < 0 {
		return "", errors.New("Invalid command for MAIL")
	}

	split := strings.Split(streamInput, ":")
	if len(split) < 2 {
		return "", errors.New("MAIL FROM command format is invalid")
	}

	from := strings.TrimSpace(strings.Join(split[1:], ""))
	this.Writer.SendOkResponse()

	return from, nil
}

/*
Function to process the RCPT TO command (constant RCPT). This command
will respond to clients with a 250 Ok response and returns an error structre
and a string containing the recipients address. Note that a client
can send one or more RCPT TO commands.
*/
func (this *SmtpWorker) Process_RCPT(streamInput string) (string, error) {
	commandCheck := strings.Index(strings.ToLower(streamInput), "rcpt to")
	if commandCheck < 0 {
		return "", errors.New("Invalid command for RCPT TO")
	}

	split := strings.Split(streamInput, ":")
	if len(split) < 2 {
		return "", errors.New("RCPT TO command format is invalid")
	}

	to := strings.TrimSpace(strings.Join(split[1:], ""))
	this.Writer.SendOkResponse()

	return to, nil
}

/*
This is the function called by the SmtpListener when a client request
is received. This will start the process by responding to the client,
start processing commands, and finally close the connection.
*/
func (this *SmtpWorker) Work() {
	go func() {
		var streamInput string
		var command smtpconstants.SmtpCommand
		var err error

		this.InitializeMailItem()
		this.Writer.SayHello()

		/*
		 * Read from the connection until we receive a QUIT command
		 * or some critical error occurs and we force quit.
		 */
		startTime := time.Now()

		for this.State != smtpconstants.SMTP_WORKER_DONE && this.State != smtpconstants.SMTP_WORKER_ERROR {
			streamInput = this.Reader.Read()
			command, err = smtpconstants.GetCommandFromString(streamInput)

			if err != nil {
				log.Println("ERROR finding command from input", streamInput, "-", err)
				this.State = smtpconstants.SMTP_WORKER_ERROR
				continue
			}

			if command == smtpconstants.QUIT {
				this.State = smtpconstants.SMTP_WORKER_DONE
				log.Println("Closing connection")
			} else {
				err = this.ExecuteCommand(command, streamInput)
				if err != nil {
					this.State = smtpconstants.SMTP_WORKER_ERROR
					log.Println("ERROR - Error executing command", command.String())
					continue
				}
			}

			if this.TimeoutHasExpired(startTime) {
				this.State = smtpconstants.SMTP_WORKER_ERROR
				continue
			}
		}

		this.Writer.SayGoodbye()
		this.Connection.Close()

		this.State = smtpconstants.SMTP_WORKER_IDLE
		this.Receiver <- this.Mail
	}()
}

/*
Determines if the time elapsed since a start time has exceeded
the command timeout.
*/
func (this *SmtpWorker) TimeoutHasExpired(startTime time.Time) bool {
	return int(time.Since(startTime).Seconds()) > smtpconstants.COMMAND_TIMEOUT_SECONDS
}
