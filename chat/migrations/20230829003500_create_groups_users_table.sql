CREATE TABLE groups_users(
    group_user_id VARCHAR(80) PRIMARY KEY,
    user_id VARCHAR(80) NOT NULL,
    group_id VARCHAR(80) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(80) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    UNIQUE(user_id, group_id)
);