# Development
## Requirements

- Go v1.13+
- Nodejs v10+
- Yarn v1.7+

## Installing Go tools

Execute these commands either:
- outside of Monitoror project
- or use `go mod tidy` after them

```shell script
# Generating mock for backend
go get -u github.com/vektra/mockery/.../

# Test utilities
go get -u gotest.tools/gotestsum

# Embed front dist into go binary
go get -u github.com/GeertJohan/go.rice/rice
```

For installing Linter, see installation guide of [golangci-lint](https://github.com/golangci/golangci-lint#install)

## Running project

Starting project:

```shell script
make install
make run
```

```shell script
cd front
yarn
yarn run serve
```

Building project:
```shell script
cd front
yarn
yarn run build
cd ..
make install
make build
```

Tests and Linting:
```shell script
make test
make lint
```

List all the available targets:
```shell script
make help
```
