// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package smtpconstants

import (
	"errors"
	"strings"
)

type SmtpCommand int

/*
Constants representing the commands that an SMTP client will
send during the course of communicating with our server.
*/
const (
	NONE SmtpCommand = iota
	DATA SmtpCommand = iota
	RCPT SmtpCommand = iota
	MAIL SmtpCommand = iota
	HELO SmtpCommand = iota
	RSET SmtpCommand = iota
	QUIT SmtpCommand = iota
)

/*
This is a command map of SMTP command strings to their int
representation. This is primarily used because there can
be more than one command to do the same things. For example,
a client can send "helo" or "ehlo" to initiate the handshake.
*/
var SmtpCommands = map[string]SmtpCommand{
	"helo":      HELO,
	"ehlo":      HELO,
	"rcpt to":   RCPT,
	"mail from": MAIL,
	"send":      MAIL,
	"rset":      RSET,
	"quit":      QUIT,
	"data":      DATA,
}

/*
Friendly string representations of commands. Useful in error
reporting.
*/
var SmtpCommandsToStrings = map[SmtpCommand]string{
	HELO: "HELO",
	RCPT: "RCPT TO",
	MAIL: "SEND",
	RSET: "RSET",
	QUIT: "QUIT",
	DATA: "DATA",
}

/*
Takes a string and returns the integer command representation. For example
if the string contains "DATA" then the value 1 (the constant DATA) will be returned.
*/
func GetCommandFromString(input string) (SmtpCommand, error) {
	result := NONE

	if input == "" {
		return result, nil
	}

	for key, value := range SmtpCommands {
		if strings.Index(strings.ToLower(input), key) > -1 {
			result = value
			break
		}
	}

	if result == NONE {
		return result, errors.New("Command " + input + " not found")
	}

	return result, nil
}

/*
Returns the string representation of a command.
*/
func (this SmtpCommand) String() string {
	return SmtpCommandsToStrings[this]
}
