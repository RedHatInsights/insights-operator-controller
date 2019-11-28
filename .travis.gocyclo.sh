#!/bin/bash

go get github.com/fzipp/gocyclo
gocyclo -over 7 -avg .
