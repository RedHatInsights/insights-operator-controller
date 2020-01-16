#!/bin/sh

# md5 for user `all` is enabled by default if POSTGRES_PASSWORD env var is present.
echo "host all postgres all md5" >> "${PGDATA}/pg_hba.conf"

echo "listen_addresses='*'" >> "${PGDATA}/postgresql.conf"
