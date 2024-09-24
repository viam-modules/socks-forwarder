#!/bin/sh
cd `dirname $0`

go build ./
exec ./socks-forwarder $@
