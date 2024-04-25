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

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kronotop/kronotop-fdb-proxy/config"
	"github.com/kronotop/kronotop-fdb-proxy/proxy"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var proxyConfig = config.Config{}

var rootCmd = &cobra.Command{
	Use:     "kronotop-fdb-proxy",
	Version: "1.0.0",
	Short:   "An MITM proxy to inspect the traffic between Kronotop and FoundationDB clusters.",
	Run: func(cmd *cobra.Command, args []string) {
		p := proxy.New(&proxyConfig)

		go func() {
			// Handle SIGINT and SIGTERM.
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			<-sigChan
			if err := p.Shutdown(); err != nil {
				log.Err(err).Msg("Failed to shutdown proxy")
			}
		}()

		if err := p.Start(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start proxy")
			if err := p.Shutdown(); err != nil {
				log.Err(err).Msg("Failed to shutdown proxy")
			}
		}
	},
}

func Execute() {
	rootCmd.Flags().DurationVarP(&proxyConfig.GracePeriod, "grace-period", "", config.DefaultGracePeriod, "maximum time period to wait before shutting down the proxy")
	rootCmd.Flags().StringVarP(&proxyConfig.Network, "network", "n", config.DefaultNetwork, "network to use")
	rootCmd.Flags().StringVarP(&proxyConfig.Host, "host", "", config.DefaultAddress, "host to bind")
	rootCmd.Flags().IntVarP(&proxyConfig.Port, "port", "p", config.DefaultPort, "port to listen")
	rootCmd.Flags().StringVarP(&proxyConfig.FdbHost, "fdb-host", "", config.DefaultFDBHost, "FDB host")
	rootCmd.Flags().IntVarP(&proxyConfig.FdbPort, "fdb-port", "", config.DefaultFDBPort, "FDB port")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run proxy")
	}
}
