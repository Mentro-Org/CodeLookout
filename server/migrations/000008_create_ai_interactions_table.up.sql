CREATE TABLE ai_interactions (
    id VARCHAR PRIMARY KEY,
    session_id VARCHAR REFERENCES ai_review_sessions(id),
    prompt VARCHAR,
    response VARCHAR,
    created_at VARCHAR
);