CREATE TABLE pull_requests (
    id VARCHAR PRIMARY KEY,
    repository_id VARCHAR REFERENCES repositories(id),
    pr_number VARCHAR,
    title VARCHAR,
    author_id VARCHAR REFERENCES users(id),
    branch VARCHAR,
    status VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);