#!/usr/bin/env bash

if getent passwd kolide &>/dev/null; then
	userdel -r kolide &>/dev/null
fi

exit 0
