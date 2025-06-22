CREATE TABLE commits (
    id VARCHAR PRIMARY KEY,
    pr_id VARCHAR REFERENCES pull_requests(id),
    commit_hash VARCHAR,
    message VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    author_id VARCHAR REFERENCES users(id)
);