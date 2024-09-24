#!/bin/sh
cd `dirname $0`

# Install .deb package (`dpkg` command should be idempotent and should
# automatically enable/start systemd service).
dpkg -i ./socks-forwarder_0.1_arm64.deb

go build ./
exec ./socks-forwarder $@
