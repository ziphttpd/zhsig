#!/bin/bash
SCRIPT=$0
TARGET_DIR=$1

if [ ! -d "${TARGET_DIR}" ]; then
	echo "Usage: ${SCRIPT} {install-directory}"
	exit 1
fi

SCRIPT_DIR=$(cd $(dirname ${SCRIPT}); pwd)
cd ${SCRIPT_DIR}
go build -o ${TARGET_DIR}/zhget cmd/zhget/zhget.go
go build -o ${TARGET_DIR}/zhsign cmd/zhsign/zhsign.go

exit 0
