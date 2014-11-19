// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package smtpio

import (
	"bytes"
	"net"
	"strings"
	"time"

	"github.com/adampresley/mailslurper/libmailslurper/smtpconstants"
)

type SmtpReader struct{
	Connection *net.TCPConn
}

/*
This function reads the raw data from the socket connection to our client. This will
read on the socket until there is nothing left to read and an error is generated.
This method blocks the socket for the number of milliseconds defined in CONN_TIMEOUT_MILLISECONDS.
It then records what has been read in that time, then blocks again until there is nothing left on
the socket to read. The final value is stored and returned as a string.
*/
func (this *SmtpReader) Read() string {
	var raw bytes.Buffer
	var bytesRead int

	bytesRead = 1

	for bytesRead > 0 {
		this.Connection.SetReadDeadline(time.Now().Add(time.Millisecond * smtpconstants.CONN_TIMEOUT_MILLISECONDS))

		buffer := make([]byte, smtpconstants.RECEIVE_BUFFER_LEN)
		bytesRead, err := this.Connection.Read(buffer)

		if err != nil {
			break
		}

		if bytesRead > 0 {
			raw.WriteString(string(buffer[:bytesRead]))
		}
	}

	return raw.String()
}

/*
This is used by the SMTP DATA command. It will read data from the connection
until the terminator is sent.
*/
func (this *SmtpReader) ReadDataBlock() string {
	var dataBuffer bytes.Buffer

	for {
		dataResponse := this.Read()

		terminatorPos := strings.Index(dataResponse, smtpconstants.SMTP_DATA_TERMINATOR)
		if terminatorPos <= -1 {
			dataBuffer.WriteString(dataResponse)
		} else {
			dataBuffer.WriteString(dataResponse[0:terminatorPos])
			break
		}
	}

	return dataBuffer.String()
}