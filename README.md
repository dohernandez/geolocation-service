# Geolocation Service

[![Coverage Status](https://coveralls.io/repos/github/dohernandez/geolocation-service/badge.svg?branch=master)](https://coveralls.io/github/dohernandez/geolocation-service?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/eb261bcd0f274be5b83b9af9f555099c)](https://www.codacy.com/app/dohernandez/geolocation-service?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dohernandez/geolocation-service&amp;utm_campaign=Badge_Grade)
[![Docker Repository on Quay](https://quay.io/repository/dohernandez/geolocation-service/status "Docker Repository on Quay")](https://quay.io/repository/dohernandez/geolocation-service)

Geolocation-service (Service) is a service responsible for displaying geolocation data. It allows to import such data and expose it via an API.

## Requirements notation
The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [RFC2119](http://tools.ietf.org/html/rfc2119).

## Table of Contents

- [Getting started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Getting the source](#getting-the-source)
    - [Development](#development)
        - [Application run](#application-run)
        - [Generate documentation](#generate-documentation)
    - [Testing](#testing)
- [Demo](#demo)
- [Troubleshooting](#troubleshooting)
    - [Known issues](#known-issues)
    
## Getting started

### Prerequisites

You need to make sure that you have `go v1.12` or later, `make` and `docker` installed

```
$ which make
/usr/bin/make
$ which go
/usr/local/bin/go
$ which docker
/usr/local/bin/docker
```

There aren't any other prerequisites needed to setup this project for development.

[[table of contents]](#table-of-contents)

### Getting the source code

Setup the project structure and fetch the repo as following:
 
```bash
go get github.com/dohernandez/geolocation-service
```

### Development

This project follows the following structure:

```markdown
|-- ci # OPTIONAL, place to store ci resources
|-- cmd # MUST be used as a main entrypoint, one folder per binary
	|-- servi # For simple cli application logic. Setup is done here
		|-- cmdimport
		        |-- cmdimport.go
	|-- servid # For simple web application logic. Setup is done here
		|-- servid.go
|-- features # OPTIONAL, place to store specification definitions in gherkin format
|-- internal # contains application specific non-reusable by any other projects code 
	|-- domain # domain packages
	|-- platform # foundational packages specific to the project
		|-- app # MUST contains base standard definitions to setup service.
		|-- http
			|-- handlers # http handler grouped by domain bundle
			|-- routes.go  # MUST contains routes
			|-- config.go # MUST contains the service configurations
			|-- container.go # MUST contains service resources
			|-- init.go # MUST initialize the service resources
		|-- storage # MUST contains the data abstraction (removing, updating, and selecting items from collection)
|-- pkg # MUST NOT import internal packages. Packages placed here should be considered as vendor.
	|-- http
		|-- rest
			|-- request
	|-- log
|-- resources # RECOMMENDED service resources. Shell helper scripts, additional files required for development, documentations.
	|-- migrations # Migration files
	|-- docs # MUST contains project documentation in human and/or machine readable format
|-- .editorconfig # OPTIONAL https://editorconfig.org
|-- .env.template # MUST contains the env variables used by the service.
|-- .travis.yml # OPTIONAL Travis CI
|-- .gitignore
|-- docker-compose.yml
|-- Dockerfile
|-- Gopkg.lock
|-- Gopkg.toml
|-- Makefile
|-- README.md
```

**Package Design**

`cmd/`
    
    * Packages that provide support for a specific program that is being built.
    * Can only import package from `internal/platform` and `pgk`.
    * Can't import package from `internal/domain`.
    * Allowed to panic an application.
    * Wrap errors with context if not being handled.
    * Majority of handling errors happen here.
    * Can recover any panic.
    * Only if system can be returned to 100% integrity.
    
`pkg`
    
    * Can't import `internal` packages. 
    * Packages placed here should be considered as vendor.
    * Stick to the testing package in go.
    * NOT allowed to panic an application.
    * NOT allowed to wrap errors.
    * Return only root cause error values.
    * NOT allowed to set policy about any application concerns.
    * NOT allowed to log. Access to trace information must be decoupled.
    * Configuration and runtime changes must be decoupled.
    * Retrieving metric and telemetry values must be decoupled.
    * Stick to the testing package in go.
    * Test files belong inside the package.
    * Focus more on unit than integration tests.
    
`internal\domain`
    
    * NOT allowed to panic an application.
    * Allowed to wrap errors when domain concern.
    * Wrap errors with context if not being handled.
    * Allowed to set policies about any application concerns.
    * Allowed to log and handle configurations natively.
    * Minority of handling errors happen here.
    * Stick to the testing package in go.
    * Test files belong inside the package.
    * Focus more on unit than integration tests.
    * Packages at the same level are not allowed to import each other.
    * Package root can import subpackages.
    * Can't import `internal\platform` package

`internal\platform`
    
    * NOT allowed to panic an application.
    * NOT allowed to set policies about application concerns.
    * NOT allowed to log. Access to trace information must be decoupled.
    * Configuration and runtime changes must be decoupled.
    * Retrieving metric and telemetry values must be decoupled.
    * Return only root cause error values.
    * Stick to the testing package in go.
    * Test files belong inside the package.
    * Focus more on unit than integration tests.
    * Packages can import each other.
    * Can import `internal\domain` package
    
This structure design is mostly inspired by [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html) by William Kennedy.

Routine operations are defined in `Makefile`.

```bash
geolocation-service routine operations

  init:                 Init the application

       -- Misc --

  build:                Build binary
  run:                  Run application (before exec this command make sure `make init` was executed)
  run-compile-daemon:   Run application with CompileDaemon (automatic rebuild on code change)
  lint:                 Check with golangci-lint
  fix-lint:             Apply goimports and gofmt
  deps:                 Ensure dependencies according to toml file
  deps-vendor:          Ensure dependencies according to lock file

       -- Environment modifiers --

  env:                  Run command with .env vars (before exec this command make sure `make init` was executed)
  envfile:              Generate .env file based on .env.template if not exists

       -- Test --

  test:                 Run tests
  test-unit:            Run unit tests
  test-integration:     Run integration tests
                       
                        Arguments:
                          TAGS     Optional tag(s) to run. Filter scenarios by tags:
                                   - "@dev": run all scenarios with wip tag
                                   - "~@notImplemented": exclude all scenarios with wip tag
                                   - "@dev && ~@notImplemented": run wip scenarios, but exclude new
                                   - "@dev,@undone": run wip or undone scenarios
                          FEATURE  Optional feature to run. Run only the specified feature.
                       
                        Examples:
                          only scenarios: 'make test-integration TAGS=@dev'
                          only one feature: 'make test-integration FEATURE=Dev'

       -- Documentation --

  docs:                 Generate api documentation (raml)

       -- Database migrations --

  create-migration:     Create database migration file, usage: "make create-migration NAME=<migration-name>"
  migrate:              Apply migrations
  migrate-cli:          Check/install migrations tool

       -- Docker --

  docker:               Run command with docker-compose (before exec this command make sure `make init` was executed)

       -- Service --

  servid-start:         Start the service (before exec this command make sure `make init` was executed)
  servid-stop:          Stop API service
  servid-log:           Display the service log

       -- CMD --

  cmdimport:            Import data from the file (before exec this command make sure `make init` was executed)
                       
                        Arguments:
                          FILE     Require file to run. Import data from the given file. Only support csv format.

Usage
  make <flags> [options]
```

[[table of contents]](#table-of-contents)

#### Application run

At first, the application MUST be initialized. Create the `.env` file with the server configuration and set up the environment variable `SERVICE_HOST_PORT`. To do so, run the command:

```bash
make init
```

```bash
>> initializing .env file
>> ensuring dependencies
```

It will create the `.env` file for you based on `.env.template` file. After the application is initialized, you can start/stop the service at any time. 

```bash
make servid-start
```

```bash
>> starting API service in port 8008 and postgres in port 5434
Creating network "geolocation-service_default" with the default driver
...
WARNING: Image for service api was built because it did not already exist. To rebuild this image you must use `docker-compose build` or `docker-compose up --build`.
Creating geolocation-service_postgres_1 ... done
Creating geolocation-service_servid_1   ... done

```

**Note** Wait a bit until the service is up and running, run `make servid-log` to check when the service is ready

```bash
servid_1    | >> running app with CompileDaemon
servid_1    | 2019/05/04 00:14:24 Running build command!
servid_1    | 2019/05/04 00:15:10 Build ok.
servid_1    | 2019/05/04 00:15:10 Restarting the given command.
servid_1    | 2019/05/04 00:15:10 stderr: {"level":"debug","message":"Creating routers","timestamp":"2019-05-04T10:15:10Z"}
servid_1    | 2019/05/04 00:15:10 stderr: {"level":"debug","message":"added `/` route","timestamp":"2019-05-04T10:15:10Z"}
servid_1    | 2019/05/04 00:15:10 stderr: {"level":"debug","message":"added `/version` route","timestamp":"2019-05-04T10:15:10Z"}
servid_1    | 2019/05/04 00:15:10 stderr: {"level":"debug","message":"added `/status` route","timestamp":"2019-05-04T10:15:10Z"}
servid_1    | 2019/05/04 00:15:10 stderr: {"level":"debug","message":"added `/health` route","timestamp":"2019-05-04T10:15:10Z"}
servid_1    | 2019/05/04 00:15:10 stderr: {"level":"debug","message":"added `/docs` route","timestamp":"2019-05-04T10:15:10Z"}
servid_1    | 2019/05/04 00:15:10 stderr: {"level":"info","message":"Starting server at port http://0.0.0.0:8000","timestamp":"2019-05-04T10:15:10Z"}
```

To stop the service run

```bash
make servid-stop
```

```bash
>> stop API service in port 8008 and postgres in port 5434
Stopping geolocation-service_servid_1 ... done
Stopping fgeolocation-service_postgres_1 ... done
Removing geolocation-service_servid_1 ... done
Removing geolocation-service_postgres_1 ... done
Removing network geolocation-service_default
```

[[table of contents]](#table-of-contents)

#### Generate documentation

Documentation items are generated using raml generator. RAML file is located `resources/raml/api.raml`. 

To update the api documentation, run `make docs`.

To see the api documentation generated, please follow the link to the api documentation under the [service root](http://localhost:8008).

```html
Welcome to geolocation-service. Please read API <a href="http://localhost:8008/docs/api.html">documentation</a>.
```


[[table of contents]](#table-of-contents)

### Testing 

Before running the test suite (unit test and behavioral test), make sure `.env` file is created and add following `docker-compose` service entries to your `/etc/hosts` (unix based systems):

```
127.0.0.1 postgres
```                                                                                                                                               

then you can run

- **Unit test**
```
make env test-unit
```

- **Integration test**
```
make docker test-integration
```

**Note** to run `Integration test` you will need to have a database running, so that why we highly recommend using `make docker <command>` to execute the suite test. Don't forget to run migration in case need it.

Another way to run the complete test suite is using docker. By using docker, there is no need to add any entries to your `/etc/hosts`:

```
make docker test
```

This is the most simple way to quickly start testing your service after cloning the repo, though it has low performance and is harder to debug.

[[table of contents]](#table-of-contents)

## Demo

Init the service if it is your first time. 

```bash
make init
```

Start the service

```bash
make servid-start
```

It will start the service on the port defined in your `.env` file. Once is up and running you can access to it thro [http://localhost:8008](http://localhost:8008).


Run migration

```bash
make docker migrate
``` 

Import geolocalation data. (data use for this example can be found in `resources/data/data_dump.csv`)

```bash
make cmdimport FILE=resources/data/data_dump.csv
...
Import statistics
time elapsed: 700.6662ms
processed: 100
accepted: 59
discarded: 41
```

then you can request the geolocation details from any of the ip defined in the file. Example:

- [http://localhost:8008/geolocation/16.70.191.240](http://localhost:8008/geolocation/16.70.191.240). You should get a json response.

```json
{
    "id": "ec4ee221-fed3-4c8e-9a73-d5aafb9c8e83",
    "ip_address": "16.70.191.240",
    "country_code": "MP",
    "country": "New Zealand",
    "city": "Port Mateo",
    "latitude": "28.208639115578364",
    "longitude": "-66.21699714924827",
    "mystery_value": 4779368757
}
```


[[table of contents]](#table-of-contents)

## Troubleshooting

### Known issues

There are no known issues.

[[table of contents]](#table-of-contents)


