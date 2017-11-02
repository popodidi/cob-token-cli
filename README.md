# cob-token-cli

> A command line tool that helps you to
> - check ETH/COB balance
> - send ETH/COB
> - allocate COBs to multiple addresses with `*.csv`

## Prerequisite

### Golang

Go: https://golang.org/doc/install

## Install

```bash
$ go get github.com/popodidi/cob-token-cli
$ go install github.com/popodidi/cob-token-cli
```

## Usage

```bash
$ cob-token-cli [command]
```

```
NAME:
   cob-token-cli - A COB token mangement command line tool

USAGE:
   cob-token-cli [global options] command [command options] [arguments...]

VERSION:
   0.1.2

COMMANDS:
     help, h  Shows a list of commands or help for one command
   private:
     send-eth      send ETHs
     send-cob      send COBs
     allocate-cob  allocate COBs to multiple addresses
   public:
     eth-balance  check ETH balance of address
     cob-balance  check COB balance of address

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

