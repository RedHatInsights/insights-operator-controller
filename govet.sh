#!/usr/bin/env bash

cd "$(dirname $0)"
go vet `go list ./...`
