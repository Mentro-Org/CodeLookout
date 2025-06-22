CREATE TABLE ai_review_sessions (
    id VARCHAR PRIMARY KEY,
    pr_id VARCHAR REFERENCES pull_requests(id),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    parameters VARCHAR
);