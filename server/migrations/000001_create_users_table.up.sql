CREATE TABLE users (
    id VARCHAR PRIMARY KEY,
    username VARCHAR,
    email VARCHAR,
    password_hash VARCHAR,
    role VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);