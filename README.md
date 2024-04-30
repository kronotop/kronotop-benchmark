# kronotop-fdb-proxy

An MITM proxy to inspect the traffic between Kronotop and FoundationDB clusters.

## Install

With a correctly configured Golang environment:

```shell
go install github.com/kronotop/kronotop-fdb-proxy@latest
```

## Usage

```
>> kronotop-fdb-proxy --help
An MITM proxy to inspect the traffic between Kronotop and FoundationDB clusters.

Usage:
  kronotop-fdb-proxy [flags]

Flags:
      --fdb-host string         FDB host (default "localhost")
      --fdb-port int            FDB port (default 4689)
      --grace-period duration   maximum time period to wait before shutting down the proxy (default 5s)
  -h, --help                    help for kronotop-fdb-proxy
      --host string             host to bind (default "localhost")
  -n, --network string          network to use (default "tcp")
  -p, --port int                port to listen (default 8080)
  -v, --version                 version for kronotop-fdb-proxy
```

## Contributions

Please don't hesitate to fork the project and send a pull request or just e-mail me to ask questions and share ideas.

## License

The Apache License, Version 2.0 - see LICENSE for more details.