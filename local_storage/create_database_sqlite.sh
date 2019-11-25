#!/bin/sh

DATABASE=controller.db

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

cat "${SCRIPT_DIR}/schema_sqlite.sql" | sqlite3 "${SCRIPT_DIR}/../${DATABASE}"
cat "${SCRIPT_DIR}/init_data_sqlite.sql" | sqlite3 "${SCRIPT_DIR}/../${DATABASE}"
