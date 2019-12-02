# Insights operator controller

[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-operator-controller)](https://goreportcard.com/report/github.com/RedHatInsights/insights-operator-controller) [![Build Status](https://travis-ci.org/RedHatInsights/insights-operator-controller.svg?branch=master)](https://travis-ci.org/RedHatInsights/insights-operator-controller)

## Description

A service to be used to store Insights operator configuration and to offer the configuration to selected operator.

## How to build the tool

Use the standard Go command:

```
go build
```

This command should create an executable file named `insights-operator-controller`.



## Start

Just run the executable file created by `go build`:

```
./insights-operator-controller
```



## Configuration

### HTTPS instead of HTTP

Change the following lines in `config.toml`:
- use_https=false
- address=":4443"

Please note that the service (when run locally) use the self-signed certificate.
You'd need to use `certs.pem` file on client side (curl, web browser etc.)


### Configuration file

Default configuration file is `config.toml`. It is possible to specify config file via environment variable
named `INSIGHTS_CONTROLLER_CONFIG_FILE`. For example:

```
export INSIGHTS_CONTROLLER_CONFIG_FILE=~/config.toml
./insights-operator-controller
```


## Data storage

Data storage used by the service is configurable via the command line parameters. Currently it is possible to configure the following data storages:
* SQLite local database: `controller.db` for the local deployment and `data.db` for functional tests
* PostgreSQL database: for local deployment and to be able to deploy the application to developer development



### SQLite

Use the following scripts from the `local_storage` subdirectory to work with SQLite database:
* `create_database_sqlite.sh` to create new database stored in file `controller.db`
* `create_test_database_sqlite.sh` to create new database stored in file `test.db`, this database will be used by tests



### PostgreSQL

PostgreSQL needs to be setup correcty:
* User `postgres` should have password set to `postgres`
* In the configuration file `/var/lib/pgsql/data/pg_hba.conf`, the method `md5` needs to be selected for user `postgres` and `all`
* The PostgreSQL daemon (service) has to be started, of course: `sudo systemctl start postgresql`

For more information how to install PostgreSQL on Fedora (or RHEL) machine, please follow this guide:
https://computingforgeeks.com/how-to-install-postgresql-on-fedora/

Use the following scripts from the `local_storage` subdirectory to work with SQLite database:
* `create_database_postgres.sh` to create new database named `controller`
* `create_test_database_postgres.sh` to create new database named `test_db`

The following two scripts can be used to drop existing database(s):
* `drop_database_postgres.sh` to drop database named `controller`
* `drop_test_database_postgres.sh` to drop database named `test_db`


### RDS AWS PostgreSQL

To set up a database on RDS AWS, an AWS account is needed. To set up the AWS account follow instructions:
https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/

After AWS account is set up, follow instructions to set up a PostreSQL database instance here : https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Tutorials.WebServerDB.CreateDBInstance.html

When database instance status becomes available, add master username, master password, and endpoint including port to `create_RDS_database_postgres.sh` script,
then run it to create database.

To drop previously created database, add the same values into `drop_RDS_database_postgres.sh` script and run it.

 
## Testing

### Unit tests

The following command run all unit tests:

```
go test ./...
```

It is also possible to increase verbosity level:

```
go test -v ./...
```

### REST API tests

REST API tests needs the running service and the test database to be prepared. In order to
perform REST API tests, start the following script:

```
./test.sh
```

Please note that the service should not be running at the same moment (as it used the same port).

## CI

[Travis CI](https://travis-ci.com/) is configured for this repository. Several tests and checks are started for all pull requests:

* Unit tests that use the standard tool `go test`
* `go fmt` tool to check code formatting. That tool is run with `-s` flag to perform [following transformations](https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
* `go vet` to report likely mistakes in source code, for example suspicious constructs, such as Printf calls whose arguments do not align with the format string.
* `golint` as a linter for all Go sources stored in this repository
* `gocyclo` to report all functions and methods with too high cyclomatic complexity. The cyclomatic complexity of a function is calculated according to the following rules: 1 is the base complexity of a function +1 for each 'if', 'for', 'case', '&&' or '||' Go Report Card warns on functions with cyclomatic complexity > 9

History of checks done by CI is available at [RedHatInsights / insights-operator-controller](https://travis-ci.org/RedHatInsights/insights-operator-controller).
