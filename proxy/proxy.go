// Copyright 2024 Kronotop
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
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
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/kronotop/kronotop-fdb-proxy/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Proxy struct {
	config   *config.Config
	listener net.Listener
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(config *config.Config) *Proxy {
	ctx, cancel := context.WithCancel(context.Background())
	return &Proxy{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Proxy) processConn(conn net.Conn) {
	defer p.wg.Done()
}

func (p *Proxy) Start() error {
	var err error

	address := net.JoinHostPort(p.config.Host, strconv.Itoa(p.config.Port))
	p.listener, err = net.Listen(p.config.Network, address)
	if err != nil {
		return errors.Wrap(err, "failed to start proxy")
	}

L:
	for {
		log.Info().Msg("Ready to accept connections on address: " + p.listener.Addr().String())
		conn, err := p.listener.Accept()
		if err != nil {
			select {
			case <-p.ctx.Done():
			default:
				log.Err(err).Msg("Failed to accept connection")
			}
		} else {
			p.wg.Add(1)
			go p.processConn(conn)
		}

		select {
		case <-p.ctx.Done():
			// Server is shutting down
			break L
		default:
		}
	}

	log.Info().Msg("Waiting for all connections to finish")

	closeCh := make(chan struct{})
	go func() {
		<-time.After(p.config.GracePeriod)
		select {
		case <-closeCh:
			return
		default:
			close(closeCh)
		}
	}()

	go func() {
		p.wg.Wait()
		select {
		case <-closeCh:
			return
		default:
			close(closeCh)
		}
	}()

	<-closeCh

	return nil
}

func (p *Proxy) Shutdown() error {
	select {
	case <-p.ctx.Done():
		return nil
	default:
	}

	log.Info().Msg("Shutting down proxy")
	p.cancel()
	err := p.listener.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close listener")
	}
	return nil
}
