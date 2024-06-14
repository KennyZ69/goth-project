-- ! PROBABLY NOT USING THIS, BUT A BASH SCRIPT, BUT WILL LEAVE IT HERE FOR A WHILE.

-- Check if the 'myapp_db' database exists and create it if not
CREATE DATABASE IF NOT EXISTS testdb;

-- Connect to the 'myapp_db' database
\c testdb;

-- Check if the 'users' table exists and create it if not
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    -- created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Check if the 'user_tokens' table exists and create it if not
CREATE TABLE IF NOT EXISTS user_tokens (
    token_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id),
    token VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);
