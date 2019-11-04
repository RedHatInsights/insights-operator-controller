#!/bin/sh

DATABASE=controller.db

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

cat "${SCRIPT_DIR}/schema.sql" | sqlite3 "${SCRIPT_DIR}/../${DATABASE}"
cat "${SCRIPT_DIR}/init_data.sql" | sqlite3 "${SCRIPT_DIR}/../${DATABASE}"
