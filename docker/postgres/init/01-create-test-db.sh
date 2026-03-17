#!/bin/sh
set -eu

if [ -z "${NOVASCANS_POSTGRES_TEST_DB:-}" ]; then
  echo "NOVASCANS_POSTGRES_TEST_DB is not set; skipping test database creation"
  exit 0
fi

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
SELECT format('CREATE DATABASE %I', '${NOVASCANS_POSTGRES_TEST_DB}')
WHERE NOT EXISTS (
  SELECT 1
  FROM pg_database
  WHERE datname = '${NOVASCANS_POSTGRES_TEST_DB}'
)\gexec
EOSQL
