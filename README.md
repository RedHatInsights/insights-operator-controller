# Insights operator controller

[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-operator-controller)](https://goreportcard.com/report/github.com/RedHatInsights/insights-operator-controller)

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



## Data storage

Data storage used by the service is configurable via the `config.tom` file. Currently it is possible to configure the following data storages:
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
