#!/bin/sh

set -e

tmpFile=$(mktemp)

( cd $(dirname "$0") &&
	go build -o "$tmpFile" app/*.go )

exec "$tmpFile" "$@"
