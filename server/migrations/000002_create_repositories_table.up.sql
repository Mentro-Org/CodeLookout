CREATE TABLE repositories (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    owner_id VARCHAR REFERENCES users(id),
    url VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);