# Makefile

SHELL:=/bin/bash
NAME=node-status-server

.PHONY: build
.PHONY: run
.PHONY: release


build:
	go build -o bin/$(NAME) -v

install:
	cp bin/$(NAME) /usr/local/bin/$(NAME)
	cp extras/$(NAME).service /etc/systemd/system/$(NAME).service

test:
	go test -v ./...

clean:
	rm -rf bin
	rm -f /etc/systemd/system/$(NAME).service

release:
	goreleaser release --snapshot --clean
