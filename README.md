# caasp-init

[![Release](https://img.shields.io/github/release/kubic-project/caasp-init.svg)](https://github.com/kubic-project/caasp-init/releases/latest)
[![Build Status](https://img.shields.io/travis/kubic-project/caasp-init/master.svg)](https://travis-ci.org/kubic-project/caasp-init)
[![codecov](https://codecov.io/gh/kubic-project/caasp-init/branch/master/graph/badge.svg)](https://codecov.io/gh/kubic-project/caasp-init)

![License: Apache 2.0](https://img.shields.io/github/license/kubic-project/caasp-init.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubic-project/caasp-init)](https://goreportcard.com/report/github.com/kubic-project/caasp-init)

Initialize docker daemon mirrors configuration.

## Usage

caasp-init will create the daemon.json configuration file and the necessary certificates for the mirror you want to config based on the kubic-init.yaml configuration file.

```shell
Usage:
  caasp-init [flags]
  caasp-init [command]

Available Commands:
  help        Help about any command
  version     Show version of caasp-init

Flags:
  -c, --config string   kubibc-init.yaml config file (default "/etc/kubic/kubic-init.yaml")
  -h, --help            help for caasp-init

Use "caasp-init [command] --help" for more information about a command.
```

## Example

`$ caasp-init`

This will use default value for configuration file

`/etc/kubic/kubic-init.yaml`

if you want to indicate a file run

`$ caasp-init -c /path/to/config/file.yaml`

If the configuration file has mirrors declared ot will generate the daemon.json
file with the following structure:

```
{
  "registries": [
    {
    "Mirrors": [
      {
      "URL": "https://airgappedregistry.com"
      }
    ],
    "Prefix": "https://mycompany.registry.com"
    }
  ],
  "iptables":false,
  "log-level": "warn"
}
```

If there is no mirror declared the configuration file will just be the default:

```
{
  "iptables":false,
  "log-level": "warn"
}
```

For help use `caasp-init help`

### help

Displays the current version of caasp-init.

### version

Displays the current version of caasp-init.

## Install

Download the latest version from [releases](https://github.com/kubic-project/caasp-init/releases/latest).

`$ tar -xvf caasp-init-linux.tgz -C /opt`

The create a symbolic link

`$ ln -s /opt/caasp-init/bin/caasp-init /usr/local/bin/caasp-init`

## Test

Clone repository into your $GOPATH. You can also use go get:

`go get github.com/kubic-project/caasp-init`

### Dependencies

* `go >= 1.11`

Note:
We use golang modules but you still need to work inside your $GOPATH for developing `caasp-init`.
Working outside GOPATH is currently **not supported**

### Running tests

To run test on this package simply run:

`make test`

#### Testing with Docker

`make test.unit`

## Code Coverage

Run first the tests. Then use `make coverage` for visualizing coverage.

Feel free to read more about this on: https://blog.golang.org/cover.

## Building

Be sure you have all prerequisites.

A simple `make` should be enough. This should compile [the main
function](cmd/root.go) and generate a `caasp-init` binary.

Your binary will be stored under `bin` folder

## Generating Releases

Run first the tests. Then use `make release` to generate the release assets.

They will be created in the [release](release/) folder.

## Community

Currently the caasp-init project lives inside the kubic echosystem.

If you have a question to ask? Want to join in the discussion? Find community information including chat and mailing lists on the main [Kubic](https://en.opensuse.org/Portal:Kubic) page.

Want to get involved but don't know what to do? Try looking at our github [issues](https://github.com/kubic-project/caasp-init/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) and use the tags `good first issue` or `help wanted`!
