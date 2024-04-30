// Copyright 2024 Kronotop
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a pipe of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxy

import (
	"net"

	"github.com/rs/zerolog/log"
)

type ConnWrapper struct {
	net.Conn
	connId int64
	label  string
}

func newConnWrapper(label string, connId int64, conn net.Conn) *ConnWrapper {
	return &ConnWrapper{
		Conn:   conn,
		connId: connId,
		label:  label,
	}
}

func (c *ConnWrapper) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	log.Info().Int64("conn_id", c.connId).Str("label", c.label).Int("read_bytes", n).Msg("")
	return n, err
}

func (c *ConnWrapper) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	log.Info().Int64("conn_id", c.connId).Str("label", c.label).Int("write_bytes", n).Msg("")
	return n, err
}
