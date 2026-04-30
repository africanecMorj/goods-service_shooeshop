#!/bin/bash
set -e

# Initialize DB if not exists
if [ ! -s /var/lib/postgresql/data/PG_VERSION ]; then
  echo "Initializing database..."
  su-exec postgres initdb -D /var/lib/postgresql/data
fi

# Start PostgreSQL in background
su-exec postgres postgres -D /var/lib/postgresql/data &

# Wait for DB
until pg_isready -h localhost -p 5432; do
  sleep 1
done

echo "PostgreSQL started"

# Run Go app (PID 1)
exec ./app