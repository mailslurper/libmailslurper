package server

import (
	"net"
	"testing"

	"github.com/mailslurper/libmailslurper/configuration"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSmtpListener(t *testing.T) {
	Convey("Creating a new SMTP Listener", t, func() {
		Convey("with a valid address and port without cert returns a listener", func() {
			var expected net.Listener
			var actual net.Listener
			var err error

			config := &configuration.Configuration{
				SmtpAddress: "127.0.0.1",
				SmtpPort:    0,
			}

			if expected, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
				panic("Unable to setup expected TCP listener: " + err.Error())
			}

			defer expected.Close()

			if actual, err = SetupSMTPServerListener(config); err != nil {
				panic("Unable to setup actual SMTP listener: " + err.Error())
			}

			So(actual.Addr().Network(), ShouldEqual, expected.Addr().Network())
			So(actual.Addr().String()[0:9], ShouldEqual, expected.Addr().String()[0:9])
		})

		Convey("with an invalid address returns an error", func() {
			var err error

			config := &configuration.Configuration{
				SmtpAddress: "abcd",
				SmtpPort:    0,
			}

			_, err = SetupSMTPServerListener(config)
			So(err.Error(), ShouldContainSubstring, "no such host")
		})

		Convey("with a valid address and port with a cert returns a listener", func() {
			var expected net.Listener
			var actual net.Listener
			var err error

			config := &configuration.Configuration{
				SmtpAddress: "127.0.0.1",
				SmtpPort:    0,
				CertFile:    "../mocks/mailslurper-cert.pem",
				KeyFile:     "../mocks/mailslurper-key.pem",
			}

			if expected, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
				panic("Unable to setup expected TCP listener: " + err.Error())
			}

			defer expected.Close()

			if actual, err = SetupSMTPServerListener(config); err != nil {
				panic("Unable to setup actual SMTP listener: " + err.Error())
			}

			So(actual.Addr().Network(), ShouldEqual, expected.Addr().Network())
			So(actual.Addr().String()[0:9], ShouldEqual, expected.Addr().String()[0:9])
		})

		Convey("with a valid address and port with an invalid cert returns an error", func() {
			var err error

			config := &configuration.Configuration{
				SmtpAddress: "127.0.0.1",
				SmtpPort:    0,
				CertFile:    "../mocks/empty-file.pem",
				KeyFile:     "../mocks/empty-file.pem",
			}

			_, err = SetupSMTPServerListener(config)
			So(err.Error(), ShouldContainSubstring, "failed to find any PEM")
		})

	})
}
