#!/usr/bin/env bash
SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
HOME_DIR="${SCRIPT_DIR}/.."
docker run "$@" -d -p 58601:8601 -v ${HOME_DIR}:/go/src/app go_test_db_worker
