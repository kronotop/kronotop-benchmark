package proxy

import (
	"net"

	"github.com/rs/zerolog/log"
)

type ConnWrapper struct {
	net.Conn
}

func newConnWrapper(conn net.Conn) *ConnWrapper {
	return &ConnWrapper{conn}
}

func (c *ConnWrapper) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	log.Info().Int("read_bytes", n).Msg(string(b[:n]))
	return n, err
}

func (c *ConnWrapper) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	log.Info().Int("write_bytes", n).Msg(string(b[:n]))
	return n, err
}
