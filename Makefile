# Makefile

SHELL:=/bin/bash
NAME=node-status-server

.PHONY: build install test clean release

build:
	go build -o bin/$(NAME)

install:
	cp bin/$(NAME) /usr/local/bin/$(NAME)
	cp extras/$(NAME).service /etc/systemd/system/$(NAME).service

test:
	go test -v ./...

clean:
	rm -rf bin
	rm -f /etc/systemd/system/$(NAME).service
	rm -rf /etc/systemd/system/$(NAME).service.d/

release:
	goreleaser release --snapshot --clean
