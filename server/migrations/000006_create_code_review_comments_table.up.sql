CREATE TABLE code_review_comments (
    id VARCHAR PRIMARY KEY,
    pr_id VARCHAR REFERENCES pull_requests(id),
    file_id VARCHAR REFERENCES files(id),
    line_number VARCHAR,
    comment_text VARCHAR,
    comment_type VARCHAR,
    created_at TIMESTAMP
);