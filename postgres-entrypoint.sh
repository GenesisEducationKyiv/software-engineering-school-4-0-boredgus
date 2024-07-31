#!/bin/bash

set -e
set -u

function create_database() {
    local database=$1

    psql -v ON_ERROR_STOP=1 \
        --username "$POSTGRES_USER" \
        -c "SELECT 'found' FROM pg_database WHERE datname = '$database';" > "$database.txt"
    if grep -q "found" "$database.txt"
    then
        echo "  Database '$database' already exists. Skipping..."
    else 
        echo "  Creating a database '$database'..."
        psql -U postgres -v ON_ERROR_STOP=1 <<-EOSQL
            CREATE DATABASE $database;
            GRANT ALL PRIVILEGES ON DATABASE $database TO $POSTGRES_USER;
            \connect $database;
            CREATE SCHEMA $database;
            ALTER SCHEMA $database OWNER TO $POSTGRES_USER;
EOSQL
    fi
}

if [ -n "$POSTGRES_DATABASES" ]; then
    echo "Multiple database creation requested: $POSTGRES_DATABASES"
    for db in $(echo $POSTGRES_DATABASES | tr ',' ' '); do
        create_database $db
    done
    echo "Multiple databases created"
fi
