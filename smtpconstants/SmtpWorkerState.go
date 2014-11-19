// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package smtpconstants

type SmtpWorkerState int

const (
	SMTP_WORKER_IDLE    SmtpWorkerState = 0
	SMTP_WORKER_WORKING SmtpWorkerState = 1
	SMTP_WORKER_DONE    SmtpWorkerState = 100
	SMTP_WORKER_ERROR   SmtpWorkerState = 101

	RECEIVE_BUFFER_LEN        = 1024
	CONN_TIMEOUT_MILLISECONDS = 5
	COMMAND_TIMEOUT_SECONDS   = 5
)
