#!/usr/bin/env bash

if ! getent passwd kolide &>/dev/null; then
	useradd --system -u 379 -U -d / -s /sbin/nologin -c kolide kolide
fi

exit 0
