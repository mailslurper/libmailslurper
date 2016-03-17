package server

import (
	"testing"

	"github.com/mailslurper/libmailslurper/customerror"
	"github.com/mailslurper/libmailslurper/mocks"
	"github.com/mailslurper/libmailslurper/model/mailitem"
	"github.com/mailslurper/libmailslurper/sanitization"
	"github.com/mailslurper/libmailslurper/smtpio"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServerPool(t *testing.T) {
	xssService := sanitization.NewXSSService()
	emailValidationService := sanitization.NewEmailValidationService()
	mockTCPConnection := mocks.NewMockTCPConn()

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
			testPool := NewServerPool(1)
			receiver := make(chan mailitem.MailItem, 10)

			/*
			 * Create a worker and set it up as if it were receiving a connection
			 */
			expected := NewSmtpWorker(1, testPool, emailValidationService, xssService)
			expected.Prepare(
				mockTCPConnection,
				receiver,
				smtpio.SmtpReader{Connection: mockTCPConnection},
				smtpio.SmtpWriter{Connection: mockTCPConnection},
			)

			/*
			 * Now ask the pool for a worker. It should perform the same steps
			 */
			actual, err := testPool.NextWorker(mockTCPConnection, receiver)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})

		Convey("No available workers in the pool returns an error", func() {
			testPool := NewServerPool(0)
			receiver := make(chan mailitem.MailItem, 10)

			expected := customerror.NoWorkerAvailable()
			_, actual := testPool.NextWorker(mockTCPConnection, receiver)

			So(actual, ShouldResemble, expected)
		})
	})
}
