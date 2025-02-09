#!/bin/bash
CURDIR=$(cd "$(dirname "$0")" || exit; pwd)
BinaryName=ActiManage.core-api
echo "$CURDIR/bin/${BinaryName}"
exec "$CURDIR/bin/${BinaryName}"
