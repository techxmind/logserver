MAKEFILE_PATH := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
GO_BIN ?= go
SHELL := /bin/bash
VERSION ?= v_$(shell git rev-parse --short HEAD)
VERSION_DATE := $(shell $(MAKEFILE_PATH)/commit_date.sh)

PWD=$(shell pwd)
OUT_DIR=$(PWD)/build

SERVICE_DIST=$(OUT_DIR)/logservice
SERVICE_SRC=service/cmd/logservice/main.go

CONSUMER_DIST=$(OUT_DIR)/logconsumer
CONSUMER_SRC=consumer/cmd/logconsumer/main.go

LDFLAGS := -s -w
LDFLAGS += -X main.version=$(VERSION)
LDFLAGS += -X main.date=$(VERSION_DATE)

.PHONY: all
all: service consumer

.PHONY: test
test:
	@$(GO_BIN) test -timeout=90s ./...

.PHONY: service
service: clean
	@echo "Build logserver"
	@$(GO_BIN) build -ldflags "$(LDFLAGS)" -o $(SERVICE_DIST) $(SERVICE_SRC)

.PHONY: consumer
consumer: clean
	@echo "Build logconsumer"
	@$(GO_BIN) build -ldflags "$(LDFLAGS)" -o $(CONSUMER_DIST) $(CONSUMER_SRC)

.PHONY: clean
clean:
	-rm -f $(OUT_DIR)/*

.PHONY: gen
gen:
	tool/gen_code.pl
