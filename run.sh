#!/bin/sh
cd `dirname $0`

exec ./socks-forwarder $@
