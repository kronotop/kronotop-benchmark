# kronotop-fdb-proxy

An MITM proxy to inspect the traffic between Kronotop and FoundationDB clusters.

## Install

With a correctly configured Golang environment:

```shell
go install github.com/kronotop/kronotop-fdb-proxy
```

## Usage

```
>> kronotop-fdb-proxy --help
A MITM proxy to inspect the traffic between Kronotop and FoundationDB clusters.

Usage:
  kronotop-fdb-proxy [flags]

Flags:
  -h, --help       help for kronotop-fdb-proxy
  -p, --port int   Port to listen (default 8080)
  -v, --version    version for kronotop-fdb-proxy
```

## Contributions

Please don't hesitate to fork the project and send a pull request or just e-mail me to ask questions and share ideas.

## License

The Apache License, Version 2.0 - see LICENSE for more details.