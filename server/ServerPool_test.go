package server

import (
	"net"
	"testing"

	"github.com/mailslurper/libmailslurper/customerror"
	"github.com/mailslurper/libmailslurper/model/mailitem"
	"github.com/mailslurper/libmailslurper/sanitization"
	"github.com/mailslurper/libmailslurper/smtpio"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServerPool(t *testing.T) {
	xssService := sanitization.NewXSSService()
	emailValidationService := sanitization.NewEmailValidationService()

	Convey("Creating a new server pool", t, func() {
		Convey("Returns a ServerPool object", func() {
			expected := make(ServerPool, 1)
			actual := NewServerPool(1)

			So(actual, ShouldHaveSameTypeAs, expected)
		})

		Convey("When created with maxWorkers 2 should have a capacity of 2", func() {
			testPool := NewServerPool(2)
			expected := 2
			actual := cap(testPool)

			So(actual, ShouldEqual, expected)
		})
	})

	Convey("Requesting a worker", t, func() {
		Convey("With available workers in the pool returns a worker", func() {
			var err error

			testPool := NewServerPool(1)
			receiver := make(chan mailitem.MailItem, 10)

			mockConnection := getMockTCPConnection()

			/*
			 * Create a worker and set it up as if it were receiving a connection
			 */
			expected := NewSmtpWorker(1, testPool, emailValidationService, xssService)
			expected.Prepare(
				mockConnection,
				receiver,
				smtpio.SmtpReader{Connection: mockConnection},
				smtpio.SmtpWriter{Connection: mockConnection},
			)

			/*
			 * Now ask the pool for a worker. It should perform the same steps
			 */
			actual, err := testPool.NextWorker(mockConnection, receiver)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})

		Convey("No available workers in the pool returns an error", func() {
			testPool := NewServerPool(0)
			receiver := make(chan mailitem.MailItem, 10)

			mockConnection := getMockTCPConnection()

			expected := customerror.NoWorkerAvailable()
			_, actual := testPool.NextWorker(mockConnection, receiver)

			So(actual, ShouldResemble, expected)
		})
	})
}

func getMockTCPConnection() net.Conn {
	var err error
	var mockConnection net.Conn
	var tcpServer net.Listener

	if tcpServer, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
		panic("Unable to create TCP listener on localhost")
	}

	go func() {
		defer tcpServer.Close()
		tcpServer.Accept()
	}()

	if mockConnection, err = net.Dial("tcp", tcpServer.Addr().String()); err != nil {
		panic("Unable to dial to mock TCP server")
	}

	return mockConnection
}
