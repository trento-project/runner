# Trento runner

The Trento runner is responsible of running the Trento configuration health checks among the installed Trento Agents. It is yet another component in the `Trento project` stack, and in order to work, it needs to talk with the Trento Web component. So, this seconds component must be available before starting using the Runner.

## Table of contents

- [Quick Start](#quick-start)
  - [Requirements](#requirements)
  - [How to run](#how-to-run)
- [Checks structure](#check-structure)
- [Development](#development)
  - [Build system](#build-system)
  - [Development dependencies](#development-dependencies)
  - [Docker](#docker)
- [Support](#support)
- [Contributing](#contributing)
- [License](#license)

## Quick Start

In order to start using the Trento Runner, follow the next instructions.

### Requirements

In order to run the Trento Runner, some requirements must be installed. Here the most important ones.
- Golang == 1.16
- Python ~=3.7
- Ansible ~=2.26.0

In order to install all the Python dependencies, run:

```shell
pip install -r requirements.dev.txt
```

### How to run

Once all the dependencies are installed and the binary is present (more instructions about how to build a development binary in the next chapters), this is how the Trento Runner works.

First of all, identify the address and port where the Trento Web component is running, as the Runner needs to communicate with him. After that, simply start the runner like:

```shell
./trento-runner start --api-host $trento-web-server --api-port $trento-web-server-port
# Run ./trento-runner -h to find additional options
```

## Checks structure

The Trento Runner configuration health checks are written in Ansible, and all of them follow a similar approach.
Find in this [documentation page](docs/runner.md) how to understand and write new checks.

## Build system

We use GNU Make as a task manager; here are some common targets:

```shell
make # clean, test and build everything

make clean # removes any build artifact
make test # executes all the tests
make fmt # fixes code formatting
make generate # refresh automatically generated code (e.g. static Go mocks)
```

Feel free to peek at the [Makefile](Makefile) to know more.

## Development dependencies

Additionally, for the development we use [`mockery`](https://github.com/vektra/mockery) for the `generate` target, which in turn is required for the `test` target.
You can install it with `go install github.com/vektra/mockery/v2`.

> Be sure to add the `mockery` binary to your `$PATH` environment variable so that `make` can find it. That usually comes with configuring `$GOPATH`, `$GOBIN`, and adding the latter to your `$PATH`.

## Docker

The [Dockerfile](Dockerfile) will automatically fetch all the required compile-time dependencies to build
the binary and finally a container image with the fresh binary in it.

You can build the component like follows:

```shell
docker build -t trento-runner .
```

# Support

Please only report bugs via [GitHub issues](https://github.com/trento-project/runner/issues);
for any other inquiry or topic use [GitHub discussion](https://github.com/trento-project/runner/discussions).

# Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

# License

Copyright 2022 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
