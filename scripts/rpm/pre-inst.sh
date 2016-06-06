#!/usr/bin/env bash

if ! getent passwd kolide &>/dev/null; then
	useradd -r -u 379 -U -d /opt/kolide -s /sbin/nologin -c kolide -m kolide
fi

exit 0
