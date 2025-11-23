INSERT INTO teams (team_name)
VALUES ('backend'),
       ('frontend');

INSERT INTO users (user_id, username, team_name, is_active)
VALUES ('user-1-alice', 'Alice', 'backend', TRUE),
       ('user-2-bob', 'Bob', 'backend', TRUE),
       ('user-3-charlie', 'Charlie', 'backend', TRUE),
       ('user-4-david', 'David', 'backend', FALSE); -- Неактивный пользователь для тестов

INSERT INTO users (user_id, username, team_name, is_active)
VALUES ('user-5-eve', 'Eve', 'frontend', TRUE),
       ('user-6-frank', 'Frank', 'frontend', TRUE),
       ('user-7-grace', 'Grace', 'frontend', TRUE);

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
VALUES
    ('pr-101', 'Feature: User Authentication', 'user-1-alice', 'OPEN', NOW() - INTERVAL '3 day', NULL),

    ('pr-102', 'Fix: Password Reset Bug', 'user-1-alice', 'MERGED', NOW() - INTERVAL '10 day', NOW() - INTERVAL '8 day'),

    ('pr-201', 'Feature: New Landing Page', 'user-5-eve', 'OPEN', NOW() - INTERVAL '1 day', NULL),

    ('pr-103', 'Refactor: Database Connection Pool', 'user-2-bob', 'OPEN', NOW() - INTERVAL '2 hour', NULL),

    ('pr-202', 'Style: Update CSS Variables', 'user-5-eve', 'MERGED', NOW() - INTERVAL '5 day', NOW() - INTERVAL '4 day');


INSERT INTO reviewers (pull_request_id, user_id)
VALUES
    ('pr-101', 'user-2-bob'),
    ('pr-101', 'user-3-charlie'),

    ('pr-102', 'user-2-bob'),
    ('pr-102', 'user-4-david'),

    ('pr-201', 'user-6-frank'),
    ('pr-201', 'user-7-grace'),

    ('pr-103', 'user-1-alice'),

    ('pr-202', 'user-6-frank');