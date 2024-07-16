#!/bin/bash

set -e
set -u

function create_user_and_database() {
    local database=$1
    echo "  Creating a database '$database'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        CREATE DATABASE $database;
        GRANT ALL PRIVILEGES ON DATABASE $database TO $POSTGRES_USER;
EOSQL
}

if [ -n "$POSTGRES_DATABASES" ]; then
    echo "Multiple database creation requested: $POSTGRE_DATABASES"
    for db in $(echo $POSTGRES_DATABASES | tr ',' ' '); do
        create_user_and_database $db
    done
    echo "Multiple databases created"
fi
