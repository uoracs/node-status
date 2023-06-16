# Makefile

SHELL:=/bin/bash
NAME="node-status-server"

export GOARCH="amd64"
export GOOS="linux"

.PHONY: build
.PHONY: run
.PHONY: release


build:
	go build -o bin/$(NAME) -v

install:
	cp bin/$(NAME) /usr/local/bin/$(NAME)
	cp extras/$(NAME).service /etc/systemd/system/$(NAME).service

clean:
	rm -rf bin
	rm -f /etc/systemd/system/$(NAME).service

release:
	goreleaser release --snapshot --clean
