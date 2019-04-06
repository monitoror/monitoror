<p align="center">
    <h2>Monitoror</h2>
</p>

<p align="center">
  <a href="https://circleci.com/gh/jsdidierlaurent/monitoror/tree/master"><img src="https://img.shields.io/circleci/project/github/jsdidierlaurent/monitoror/master.svg?style=for-the-badge" alt="Build"/></a>
  <a href="https://codecov.io/gh/jsdidierlaurent/monitoror"><img src="https://img.shields.io/codecov/c/gh/jsdidierlaurent/monitoror/master.svg?style=for-the-badge" alt="Code Coverage"/></a>
  <a href="https://github.com/jsdidierlaurent/monitoror/releases"><img src="https://img.shields.io/github/release/jsdidierlaurent/monitoror.svg?style=for-the-badge" alt="Releases"/></a>
  <br>
  <img src="https://img.shields.io/badge/Go-1.12-blue.svg?style=for-the-badge" alt="Go"/>
  <img src="https://img.shields.io/badge/NodeJS-10.0-blue.svg?style=for-the-badge" alt="NodeJS"/>
</p>

------------------------------------

## Introduction

TODO

## Documentation

## Contribution

## Development
### Requirement
- Go v1.12+
- Nodejs v10+
- Yarn v1.7+
- GNU make

### Installing tools
```bash
# Go tools, execute this command outside of monitoror project or use `go mod tidy` after
go get github.com/vektra/mockery          # Generating mock for backend
go get gotest.tools/gotestsum             # Tests utilities
go get github.com/GeertJohan/go.rice/rice # Pack front dist into go source

# Node tools
yarn global add @vue/cli  # Not necessarily required
```

### Running project
```bash
# Front
cd front
yarn install
yarn run serve
```

```bash
# Back
make install
make run
# Or
make run-faker
```

### Building project
```bash
cd front
yarn install
yarn run build
cd ..
make install
make build
```

### Generating mocks
```bash
# For generating monitorable mocks
make mock

# For generating all mocks (only needed if golang net interface change)
make mock-all
```

## License
This project is under [MIT license](LICENSE).
