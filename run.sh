#!/bin/bash
#
# Runs the web server using production settings.
#
set -euo pipefail

export HTTPS_CERT_FILE=/etc/ssl/private/hkjn.me.crt
export HTTPS_KEY_FILE=/etc/ssl/private/hkjn.me.key
export BIND_ADDR=:4430
export PROD=1
go get -u hkjn.me/hkjnweb/...
go build ./cmd/server/hkjnserver.go
if pgrep hkjnserver 1>/dev/null; then
	kill $(pgrep hkjnserver)
fi
./hkjnserver -alsologtostderr
