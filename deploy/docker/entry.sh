#!/bin/sh

/usr/sbin/sshd -D -e &

evm-lite "$@"

