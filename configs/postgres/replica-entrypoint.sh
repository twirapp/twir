#!/bin/bash
set -e

DATA_DIR="/var/lib/postgresql"
chown -R postgres:postgres "$DATA_DIR"

echo "postgres:5432:*:replicator:replicator" > /var/lib/postgresql/.pgpass

if [ -z "$(ls -A $DATA_DIR)" ]; then
  echo "Initializing replica from primary..."

  until pg_isready -h postgres -p 5432; do
    sleep 2
  done

  echo "Primary is available. Starting base backup..."

  gosu postgres -c "pg_basebackup \
      -h postgres \
      -D $DATA_DIR \
      -U replicator \
      -Fp \
      -Xs \
      -P \
      -R"

  echo "Replica initialized"
fi

exec gosu postgres postgres \
  -c log_statement=all \
  -c log_connections=on
