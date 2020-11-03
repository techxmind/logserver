GO_BIN ?= go
SHELL := /bin/bash

PWD=$(shell pwd)
OUT_DIR=$(PWD)/build

SERVICE_DIST=$(OUT_DIR)/logservice
SERVICE_SRC=service/cmd/logservice/main.go

CONSUMER_DIST=$(OUT_DIR)/logconsumer
CONSUMER_SRC=consumer/cmd/logconsumer/main.go

.PHONY: all
all: service consumer

.PHONY: test
test:
	@$(GO_BIN) test -timeout=90s ./...

.PHONY: service
service: clean test
	@echo "Build logserver"
	@$(GO_BIN) build -o $(SERVICE_DIST) $(SERVICE_SRC)

.PHONY: consumer
consumer: clean test
	@echo "Build logconsumer"
	@$(GO_BIN) build -o $(CONSUMER_DIST) $(CONSUMER_SRC)

.PHONY: clean
clean:
	-rm -f $(OUT_DIR)/*
