<p align="center">
    <h2>Monitoror</h2>
</p>

<p align="center">
  <a href="https://travis-ci.org/monitoror/monitoror/branches"><img src="https://img.shields.io/travis/monitoror/monitoror/master.svg?style=for-the-badge" alt="Build"/></a>
  <a href="https://codecov.io/gh/monitoror/monitoror"><img src="https://img.shields.io/codecov/c/gh/monitoror/monitoror/master.svg?style=for-the-badge" alt="Code Coverage"/></a>
  <a href="https://github.com/monitoror/monitoror/releases"><img src="https://img.shields.io/github/release/monitoror/monitoror.svg?style=for-the-badge" alt="Releases"/></a>
  <br>
  <img src="https://img.shields.io/badge/Go-1.12-blue.svg?style=for-the-badge" alt="Go"/>
  <img src="https://img.shields.io/badge/NodeJS-10.0-blue.svg?style=for-the-badge" alt="NodeJS"/>
</p>

------------------------------------

TODO: Introduction, Documentation, Contributing



## Development

### Requirements

- Go v1.12+
- Nodejs v10+
- Yarn v1.7+
- GNU make


### Installing Go tools

Execute these commands either:
- outside of Monitoror project
- or use `go mod tidy` after them

```bash
# Generating mock for backend
go get github.com/vektra/mockery/.../
# Test utilities
go get gotest.tools/gotestsum
# Embed front dist into go binary
go get github.com/GeertJohan/go.rice/rice
```


### Running project
```bash
# Front
cd front
yarn
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
yarn
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
