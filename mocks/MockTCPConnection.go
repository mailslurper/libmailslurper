package mocks

import (
	"errors"
	"io"
	"net"
	"time"
)

type MockTCPConn struct {
	Reader io.Reader
	Writer io.Writer
	Closer io.Closer

	ReadShouldError             bool
	WriteShouldError            bool
	CloseShouldError            bool
	LocalAddress                net.Addr
	RemoteAddress               net.Addr
	SetDeadlineShouldError      bool
	SetReadDeadlineShouldError  bool
	SetWriteDeadlineShouldError bool
}

func NewMockTCPConn() *MockTCPConn {
	return &MockTCPConn{
		ReadShouldError:             false,
		WriteShouldError:            false,
		CloseShouldError:            false,
		LocalAddress:                NewMockIPAddress("localhost", "127.0.0.1"),
		RemoteAddress:               NewMockIPAddress("localhost", "127.0.0.1"),
		SetDeadlineShouldError:      false,
		SetReadDeadlineShouldError:  false,
		SetWriteDeadlineShouldError: false,
	}
}

func (mockTCPConn *MockTCPConn) Read(b []byte) (int, error) {
	if mockTCPConn.ReadShouldError {
		return 0, errors.New("Read error")
	}

	return 1, nil
}

func (mockTCPConn *MockTCPConn) Write(b []byte) (int, error) {
	if mockTCPConn.WriteShouldError {
		return 0, errors.New("Write error")
	}

	return 1, nil
}

func (mockTCPConn *MockTCPConn) Close() error {
	if mockTCPConn.CloseShouldError {
		return errors.New("Close failed")
	}

	return nil
}

func (mockTCPConn *MockTCPConn) LocalAddr() net.Addr {
	return mockTCPConn.LocalAddress
}

func (mockTCPConn *MockTCPConn) RemoteAddr() net.Addr {
	return mockTCPConn.RemoteAddress
}

func (mockTCPConn *MockTCPConn) SetDeadline(t time.Time) error {
	if mockTCPConn.SetDeadlineShouldError {
		return errors.New("SetDeadline error")
	}

	return nil
}

func (mockTCPConn *MockTCPConn) SetReadDeadline(t time.Time) error {
	if mockTCPConn.SetReadDeadlineShouldError {
		return errors.New("SetReadDeadline error")
	}

	return nil
}

func (mockTCPConn *MockTCPConn) SetWriteDeadline(t time.Time) error {
	if mockTCPConn.SetWriteDeadlineShouldError {
		return errors.New("SetWriteDeadline error")
	}

	return nil
}
