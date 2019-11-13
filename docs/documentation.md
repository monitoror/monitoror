# Installation

Working and tested on Linux, macOS, Windows and RaspberryPi3.

**Note**: Doesn't require Go.

There is no installation script yet, but you can download latest [release](https://github.com/monitoror/monitoror/releases/latest) manually.

If you want to install it on RaspberryPi and auto start Monitoror on boot, you can follow [docs/installation/raspberrypi.md](./installation/raspberrypi.md)

## Configuration

There is two types of configuration. 
* Core configuration: credentials, timeouts, cache ...
* Ui configuration: tiles and layout definition

### Core configuration

Backend configuration is set with **environment variables**. You can also use **.env** file. beside of monitoror binary.

| Env | Default | Description |
|-----|---------|-------------|
| MO_PORT | 8080       | application port |
| MO_ENV | production | |
| MO_UPSTREAMCACHEEXPIRATION | 10000 | UpstreamCache is used to respond before executing the request. Avoid overloading services. |
| MO_DOWNSTREAMCACHEEXPIRATION | 120000 | DownstreamCache is used to respond after executing the request in case of timeout. Avoid false negative on CI check. |

For every monitored services, you have env, but they are defined in ihe section [Tiles Definition](#Tiles-Definition)

### UI configuration

This configuration must be stored in .json file and passed to the front.

There is two way to give the configuration to the front :
* By his path: http://localhost:8080?configPath=./config.json
* By url (if you host you file): http://localhost:8080?configUrl=https://gist.githubusercontent.com/monitoror/xxx/raw/xxx/config.json

The file is structured like this:
```json
{
  "columns": 4,
  "tiles": [
    { "type": "EMPTY", "columnSpan": 2, "rowSpan": 2},
    { "type": "PING", "label": "Awesome ping" ,"params": { "hostname": "127.0.0.1" } },
    { "type": "GROUP", "tiles": [
      { "type": "PING", "label": "Awesome ping inside group" ,"params": { "hostname": "127.0.0.1" } }
    ]}
  ]
}
```

Each tile is defined by his type. You can find every monitorable tile type in the section [Tiles Definition](#Tiles-Definition)

## Tiles Definition
### Common

Each tile definition can override some layout information:
* label: default defined by the backend (exemple: for PING, it's the hostname)
* columnSpan: default 1
* rowSpan: default 1

### Empty
```json
{ "type": "EMPTY"}
```

### Group
```json
{ "type": "GROUP", "tiles": [
  ...
]}
```

### Ping
Core configuration:

| Env | Default | Description |
|-----|---------|-------------|
| MO_MONITORABLE_PING_COUNT | 2 | number of ping send |
| MO_MONITORABLE_PING_TIMEOUT | 1000 | timeout before returning error |
| MO_MONITORABLE_PING_INTERVAL | 100 | interval between ping |

Ui configuration:

| Params | Default | Optional |
|--------|---------|----------|
| hostname | | false |

Exemple :

```json
{ "type": "PING", "params": { "hostname": "localhost" } }
```

### Port
Core configuration:

| Env | Default | Description |
|-----|---------|-------------|
| MO_MONITORABLE_PORT_TIMEOUT | 1000 | timeout before returning error |

Ui configuration:

| Params | Default | Optional |
|--------|---------|----------|
| hostname | | false |
| port | | false |

Exemple :
```json
{ "type": "PORT", "params": { "hostname": "localhost", "port": 1234 } }
```

### HTTP
Core configuration:

| Env | Default | Description |
|-----|---------|-------------|
| MO_MONITORABLE_HTTP_TIMEOUT | 1000 | timeout before returning error |
| MO_MONITORABLE_HTTP_SSLVERIFY | true | check if ssl certificate is valid |

#### HTTP-ANY

Ui configuration:

| Params | Default | Optional |
|--------|---------|----------|
| url | | false |
| statusCodeMin | 200 | true |
| statusCodeMax | 399 | true |

Exemple :

```json
{ "type": "HTTP-ANY", "params": { "url": "http://localhost/test", "statusCodeMin": 200, "statusCodeMax": 299 } }
```

#### HTTP-RAW

Ui configuration:

| Params | Default | Optional |
|--------|---------|----------|
| url | | false |
| statusCodeMin | 200 | true |
| statusCodeMax | 399 | true |
| regex | | true | 

Exemple :

```json
{ "type": "HTTP-RAW", "params": { "url": "http://localhost/test", "regex": ".*(" } }
```

#### HTTP-JSON

Ui configuration:

| Params | Default | Optional |
|--------|---------|----------|
| url | | false |
| statusCodeMin | 200 | true |
| statusCodeMax | 399 | true |
| key | | false |
| regex | | true | 

Exemple :

```json
{ "type": "HTTP-JSON", "params": { "url": "http://localhost/test", "key":".bloc1.\"bloc.2\".[0].value2", "regex": ".*(" } }
```

### Pingdom

Core configuration:

| Env | Default | Description |
|-----|---------|-------------|
| MO_MONITORABLE_PINGDOM_URL | https://api.pingdom.com/api/3.1 | pingdom api base url |
| MO_MONITORABLE_PINGDOM_TOKEN | | your private api token |
| MO_MONITORABLE_PINGDOM_TIMEOUT | 1000 | timeout before returning error |
| MO_MONITORABLE_PINGDOM_CACHEEXPIRATION | 30000 | specific cache duration for pingdom check |

#### PINGDOM-CHECK

Ui configuration:

| Params | Default | Optional |
|--------|---------|----------|
| id | | false |

Exemple :

```json
{ "type": "PINGDOM-CHECK", "params": { "id": 10 } }
```

#### PINGDOM-CHECKS

**Dynamic Tile**

Ui configuration:

| Params | Default | Optional | Comment |
|--------|---------|----------|---------|
| tags | | true | |
| sortBy | | true | only support empty or "name" |

Exemple :

```json
{ "type": "PINGDOM-CHECKS", "params": { "tags": "eu-west", "sortBy": "name" } }
```

### Jenkins

Core configuration:

| Env | Default | Description |
|-----|---------|-------------|
| MO_MONITORABLE_JENKINS_URL | | jenkins base url |
| MO_MONITORABLE_JENKINS_LOGIN | | your login |
| MO_MONITORABLE_JENKINS_TOKEN | | your private api token |
| MO_MONITORABLE_JENKINS_TIMEOUT | 2000 | timeout before returning error |
| MO_MONITORABLE_JENKINS_SSLVERIFY | true | specific cache duration for pingdom check |

#### JENKINS-BUILD

Ui configuration:

| Params | Default | Optional | Comment |
|--------|---------|----------|---------|
| job | | false | |
| branch | | true | |

Exemple :

```json
{ "type": "JENKINS-BUILD", "params": { "job": "test-job", "branch": "master" } }
```

#### JENKINS-MULTIBRANCH

**Dynamic Tile**

Ui configuration:

| Params | Default | Optional | Comment |
|--------|---------|----------|---------|
| job | | false | |
| match | | true | |
| unmatch | | true | |

Exemple :

```json
{ "type": "JENKINS-MULTIBRANCH", "params": { "job": "eu-west", "match": "feat/*" } }
```

### TravisCI

Core configuration:

| Env | Default | Description |
|-----|---------|-------------|
| MO_MONITORABLE_TRAVISCI_URL | https://api.travis-ci.org/ | travisci api base url |
| MO_MONITORABLE_TRAVISCI_TOKEN | | your private api token |
| MO_MONITORABLE_TRAVISCI_TIMEOUT | 2000 | timeout before returning error |

Ui configuration:

| Params | Default | Optional | Comment |
|--------|---------|----------|---------|
| group | | false | |
| repository | | false | |
| branch | | false | |

Exemple :

```json
{ "type": "TRAVISCI-BUILD", "params": { "group": "group", "repository": "test", "branch": "master" } }
```

### Azure DevOps

Core configuration:

| Env | Default | Description |
|-----|---------|-------------|
| MO_MONITORABLE_AZUREDEVOPS_URL | | azure devOps base url |
| MO_MONITORABLE_AZUREDEVOPS_TOKEN | | your private api token |
| MO_MONITORABLE_AZUREDEVOPS_TIMEOUT | 4000 | timeout before returning error |

#### AZUREDEVOPS-BUILD

Ui configuration:

| Params | Default | Optional | Comment |
|--------|---------|----------|---------|
| project | | false | |
| definition | | false | |
| branch | | true | |

Exemple :

```json
{ "type": "AZUREDEVOPS-BUILD", "params": { "project": "project", "definition": 1, "branch": "master" } }
```

#### AZUREDEVOPS-RELEASE

Ui configuration:

| Params | Default | Optional | Comment |
|--------|---------|----------|---------|
| project | | false | |
| definition | | false | |

Exemple :

```json
{ "type": "AZUREDEVOPS-RELEASE", "params": { "project": "project", "definition": 1 } }
```
