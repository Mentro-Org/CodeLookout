CREATE TABLE files (
    id VARCHAR PRIMARY KEY,
    pr_id VARCHAR REFERENCES pull_requests(id),
    path VARCHAR,
    status VARCHAR
);