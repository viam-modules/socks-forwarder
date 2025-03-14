#!/bin/sh
cd `dirname $0`

echo "Installing the SOCKS forwarder systemd service if it has not already been installed"
sudo dpkg -E -i ./socks-forwarder_0.5_arm64.deb

exec ./socks-forwarder-module $@
