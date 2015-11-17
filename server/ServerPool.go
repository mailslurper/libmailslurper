// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package server

import (
	"log"
	"net"

	"github.com/mailslurper/libmailslurper/model/mailitem"
	"github.com/mailslurper/libmailslurper/sanitization"
	"github.com/mailslurper/libmailslurper/smtpio"
)

/*
ServerPool represents a pool of SMTP workers. This will
manage how many workers may respond to SMTP client requests
and allocation of those workers.
*/
type ServerPool chan *SmtpWorker

/*
JoinQueue adds a worker to the queue.
*/
func (pool ServerPool) JoinQueue(worker *SmtpWorker) {
	pool <- worker
}

/*
Create a new server pool with a maximum number of SMTP
workers. An array of workers is initialized with an ID
and an initial state of SMTP_WORKER_IDLE.
*/
func NewServerPool(maxWorkers int) ServerPool {
	xssService := sanitization.NewXSSService()
	emailValidationService := sanitization.NewEmailValidationService()

	pool := make(ServerPool, maxWorkers)

	for index := 0; index < maxWorkers; index++ {
		pool.JoinQueue(NewSmtpWorker(
			index+1,
			pool,
			emailValidationService,
			xssService,
		))
	}

	log.Println("libmailslurper: INFO - Worker pool configured for", maxWorkers, "worker(s)")
	return pool
}

/*
NextWorker retrieves the next available worker from
the queue.
*/
func (pool ServerPool) NextWorker(connection *net.TCPConn, receiver chan mailitem.MailItem) *SmtpWorker {
	/*
	 * TODO: This blocks until a worker is available. Perhaps implement a timeout?
	 */
	worker := <-pool
	worker.Prepare(
		connection,
		receiver,
		smtpio.SmtpReader{Connection: connection},
		smtpio.SmtpWriter{Connection: connection},
	)

	log.Println("libmailslurper: INFO - Worker", worker.WorkerId, "queued to handle connection from", connection.RemoteAddr().String())
	return worker
}
