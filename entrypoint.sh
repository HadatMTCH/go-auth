#!/bin/sh
set -e

# Wait for postgres
until nc -z $DB_HOST 5432; do
    echo "Waiting for postgres..."
    sleep 1
done

# Run migrations
goose -dir migrations postgres "host=$DB_HOST user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable" up

# Start application
exec "$@"