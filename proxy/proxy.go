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
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/kronotop/kronotop-fdb-proxy/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Proxy represents a proxy server that forwards network traffic between Kronotop and FoundationDB clusters.
type Proxy struct {
	config   *config.Config
	listener net.Listener
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// New is a function that creates a new instance of the Proxy struct.
func New(config *config.Config) *Proxy {
	ctx, cancel := context.WithCancel(context.Background())
	return &Proxy{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

// processConn is a method of the Proxy struct that is responsible for handling a connection.
func (p *Proxy) processConn(conn net.Conn) {
	defer p.wg.Done()

	addr := net.JoinHostPort(p.config.FdbHost, strconv.Itoa(p.config.FdbPort))
	serverConn, err := net.Dial(p.config.Network, addr)
	if err != nil {
		panic(err)
	}

	clientConn := newConnWrapper(conn)

	go func() {
		n, err := io.Copy(serverConn, clientConn)
		if err != nil {
			log.Err(err).Msg("Failed to accept connection")
			return
		}
		log.Info().Int64("copied bytes", n)
	}()

	go func() {
		n, err := io.Copy(clientConn, serverConn)
		if err != nil {
			log.Err(err).Msg("Failed to accept connection")
			return
		}
		log.Info().Int64("copied bytes", n)
	}()

}

// discoverHostAddress is a method of the Proxy struct that is responsible for discovering the host address.
// It checks if a specific interface is defined in the Proxy's config structure. If so, it retrieves the addresses
// of that interface and searches for a suitable host address. The first non-link local unicast IP address found is
// assigned to the Proxy's config.Host variable. If no suitable address is found and p.config.Host is still empty,
// it returns an error indicating that no host address is specified.
func (p *Proxy) discoverHostAddress() error {
	if p.config.Interface != "" {
		log.Info().Str("interface", p.config.Interface).Msg("Discovering interface addresses")
		if iface, _ := net.InterfaceByName("en0"); iface != nil {
			addresses, _ := iface.Addrs()
			for _, address := range addresses {
				ipNet, ok := address.(*net.IPNet)
				if ok {
					if ipNet.IP.IsLinkLocalUnicast() {
						continue
					}
					p.config.Host = ipNet.IP.String()
					break
				}
			}
		}
	}

	if p.config.Host == "" {
		return errors.New("no host specified")
	}
	return nil
}

// Start is a method of the Proxy struct that is responsible for starting the proxy server.
func (p *Proxy) Start() error {
	err := p.discoverHostAddress()
	if err != nil {
		return errors.Wrap(err, "failed to start proxy")
	}

	address := net.JoinHostPort(p.config.Host, strconv.Itoa(p.config.FdbPort))
	p.listener, err = net.Listen(p.config.Network, address)
	if err != nil {
		return errors.Wrap(err, "failed to start proxy")
	}

	log.Info().Msg("Ready to accept connections on address: " + p.listener.Addr().String())
L:
	for {
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

// Shutdown is a method of the Proxy struct that is responsible for gracefully shutting down the proxy.
// It first checks if the context for the proxy is done. If so, it returns immediately.
// Otherwise, it logs a message indicating that the proxy is shutting down, cancels the context, and closes the listener.
// If there is an error closing the listener, it is wrapped in an error and returned.
// The method returns nil if the shutdown was successful.
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
