# Insights operator controller

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
