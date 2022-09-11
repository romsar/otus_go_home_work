package main

import (
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
	addr    string
	timeout time.Duration
}

func (tc *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.addr, tc.timeout)
	if err != nil {
		return errors.Wrap(err, "go-telnet connect")
	}

	tc.conn = conn

	return nil
}

func (tc *telnetClient) Close() error {
	if err := tc.conn.Close(); err != nil {
		return errors.Wrap(err, "go-telnet close connection")
	}

	return nil
}

func (tc *telnetClient) Send() error {
	if _, err := io.Copy(tc.conn, tc.in); err != nil {
		return errors.Wrap(err, "go-telnet send")
	}

	return nil
}

func (tc *telnetClient) Receive() error {
	if _, err := io.Copy(tc.out, tc.conn); err != nil {
		return errors.Wrap(err, "go-telnet receive")
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		in:      in,
		out:     out,
		addr:    address,
		timeout: timeout,
	}
}
