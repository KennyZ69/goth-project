#!/bin/bash

# ?? Remember to make your script executable by running chmod +x setup_postgres.sh

source ./.env
export PGPASSWORD=$PSQ_PWD

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
    username VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_details (
    user_id INT REFERENCES users(user_id),
    bio TEXT,
    profile_image TEXT,
    role VARCHAR(100),
    country VARCHAR(100),
    city VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS user_tokens (
    token_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id),
    token VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS connection_requests (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    status VARCHAR(50) NOT NULL, -- "pending", "accepted", "denied"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(user_id),
    FOREIGN KEY (receiver_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS friends (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    user_name VARCHAR(50) NOT NULL,
    friend_id INT NOT NULL,
    friend_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (friend_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS chats (
    chat_id SERIAL PRIMARY KEY,
    user1_id INT REFERENCES users(user_id),
    user2_id INT REFERENCES users(user_id),
    last_message TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    chat_id INT REFERENCES chats(chat_id),
    sender_id INT REFERENCES users(user_id),
    receiver_id INT REFERENCES users(user_id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

if ! table_exists "user_details"; then
    echo "Table 'user_details' does not exist. Creating table."
    create_tables
fi

if ! table_exists "user_tokens"; then
    echo "Table 'user_tokens' does not exist. Creating table."
    create_tables
fi

echo "Database and tables are set up."
