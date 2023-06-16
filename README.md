# node-status

Tiny HTTP server that responds with node status information

## Usage

A single endpoint that returns status information based on some ansible facts.

Example:

```bash
$ curl localhost:8080/
{"name":"myhost","ansible_status":"ok","provision_status":"success"}
```

This service probably isn't useful to anyone but us :smile:

## Installation

Clone the repo, `make`, `make test`, `make install`. You'll also need to `systemctl daemon-reload` and `systemctl start node-status-server`.

## Configuration

It's not really configurable outside the hosting parameters. By default, it listens on `0.0.0.0:8080`.

If you'd like to change the port, you can run `systemctl edit node-status-server`, which creates a file at `/etc/systemd/system/node-status-server.d/override.conf`, with the following content:

```
[Service]
Environment="NODE_STATUS_SERVER_HOST=myhost"
Environment="NODE_STATUS_SERVER_PORT=9000"
```
