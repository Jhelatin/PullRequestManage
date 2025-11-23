CREATE TABLE IF NOT EXISTS teams
(
    team_name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users
(
    user_id   VARCHAR(255) PRIMARY KEY,
    username  VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) REFERENCES teams (team_name) ON DELETE SET NULL,
    is_active BOOLEAN      NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pull_requests
(
    pull_request_id   VARCHAR(255) PRIMARY KEY,
    pull_request_name TEXT         NOT NULL,
    author_id         VARCHAR(255) REFERENCES users (user_id),
    status            VARCHAR(50)  NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    merged_at         TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS reviewers
(
    pull_request_id VARCHAR(255) REFERENCES pull_requests (pull_request_id) ON DELETE CASCADE,
    user_id         VARCHAR(255) REFERENCES users (user_id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, user_id)
);

CREATE INDEX idx_reviewers_user_id ON reviewers(user_id);