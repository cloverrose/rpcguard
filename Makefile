export

# Setup Go Variables
GOPATH := $(shell go env GOPATH)
GOBIN := $(PWD)/bin

# Invoke shell with new path to enable access to bin
PATH := $(GOBIN):$(PATH)
PATH := $(shell aqua root-dir)/bin:$(PATH)
SHELL := env "PATH=$(PATH)" bash

# aqua/install installs aqua dependencies.
.PHONY: aqua/install
aqua/install:
	aqua install

# aqua/update-checksum updates the checksums for the aqua dependencies.
.PHONY: aqua/update-checksum
aqua/update-checksum:
	aqua update-checksum --deep --prune

# aqua/reset removes all the aqua dependencies.
.PHONY: aqua/reset
aqua/reset:
	aqua rm -all
