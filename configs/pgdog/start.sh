#!/bin/sh

set -e  # Exit on any error

# Read secrets
POSTGRES_USER=$(cat /run/secrets/twir_postgres_user)
POSTGRES_PASSWORD=$(cat /run/secrets/twir_postgres_password)
POSTGRES_DB=$(cat /run/secrets/twir_postgres_db)

PGDOG_USERS_BASE="/pgdog/users_base.toml"
PGDOG_USERS_FILE="/pgdog/users.toml"

echo "Applying secrets to $PGDOG_USERS_BASE"
echo "DB: $POSTGRES_DB, User: $POSTGRES_USER"

cp "$PGDOG_USERS_BASE" "$PGDOG_USERS_FILE"

# Use Linux-compatible sed -i (no backup suffix)
sed -i "s/database = \".*\"/database = \"$POSTGRES_DB\"/" "$PGDOG_USERS_FILE"
echo "Database sed applied"
sed -i "s/password = \".*\"/password = \"$POSTGRES_PASSWORD\"/" "$PGDOG_USERS_FILE"

echo "Config updated. Starting pgdog..."

exec /usr/local/bin/pgdog
