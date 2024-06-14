#!/bin/bash

# ?? Remember to make your script executable by running chmod +x setup_postgres.sh

source ./.env
export PSQ_PWD

# Variables
DB_NAME="testdb"
DB_USER="postgres" # or your custom user

# Function to check if a database exists
database_exists() {
    psql -U "$DB_USER" -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1
}

# Function to create database if it doesn't exist
create_database() {
    psql -U "$DB_USER" -c "CREATE DATABASE $DB_NAME"
}

# Function to check if a table exists
table_exists() {
    psql -U "$DB_USER" -d "$DB_NAME" -tc "SELECT 1 FROM pg_tables WHERE tablename = '$1'" | grep -q 1
}

# Function to create tables if they don't exist
create_tables() {
    psql -U "$DB_USER" -d "$DB_NAME" <<EOF
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tokens (
    token_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id),
    token VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);
EOF
}

# Main script execution
if ! database_exists; then
    echo "Database does not exist. Creating database: $DB_NAME"
    create_database
fi

if ! table_exists "users"; then
    echo "Table 'users' does not exist. Creating table."
    create_tables
fi

if ! table_exists "user_tokens"; then
    echo "Table 'user_tokens' does not exist. Creating table."
    create_tables
fi

echo "Database and tables are set up."