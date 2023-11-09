CREATE TABLE group_messages(
   message_id VARCHAR(80) PRIMARY KEY,
   group_id VARCHAR(80) NOT NULL,
   user_id VARCHAR(80) NOT NULL,
   content TEXT NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   deleted_at TIMESTAMPTZ NULL
);