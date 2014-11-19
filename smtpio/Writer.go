// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package smtpio

import (
	"log"
	"net"

	"github.com/adampresley/mailslurper/libmailslurper/smtpconstants"
)

type SmtpWriter struct{
	Connection *net.TCPConn
}

/*
Function to tell a client that we are done communicating. This sends
a 221 response. It returns true/false for success and a string
with any response.
*/
func (this *SmtpWriter) SayGoodbye() error {
	return this.SendResponse(smtpconstants.SMTP_CLOSING_MESSAGE)
}

/*
Sends a hello message to a new client. The SMTP protocol
dictates that you must be polite. :)
*/
func (this *SmtpWriter) SayHello() error {
	err := this.SendResponse(smtpconstants.SMTP_WELCOME_MESSAGE)
	if err != nil {
		return err
	}

	log.Println("Reading data from client connection...")
	return nil
}

func (this *SmtpWriter) SendDataResponse() error {
	return this.SendResponse(smtpconstants.SMTP_DATA_RESPONSE_MESSAGE)
}

/*
Function to send a response to a client connection. It returns true/false for success and a string
with any response.
*/
func (this *SmtpWriter) SendResponse(response string) error {
	_, err := this.Connection.Write([]byte(string(response + smtpconstants.SMTP_CRLF)))
	return err
}

func (this *SmtpWriter) SendHELOResponse() error {
	return this.SendResponse(smtpconstants.SMTP_HELLO_RESPONSE_MESSAGE)
}

func (this *SmtpWriter) SendOkResponse() error {
	return this.SendResponse(smtpconstants.SMTP_OK_MESSAGE)
}
