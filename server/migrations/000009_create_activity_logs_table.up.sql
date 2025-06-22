CREATE TABLE activity_logs (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR REFERENCES users(id),
    event_type VARCHAR,
    metadata VARCHAR,
    timestamp VARCHAR
);