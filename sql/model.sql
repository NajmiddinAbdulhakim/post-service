CREATE TABLE IF NOT EXISTS posts (
    id uuid NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS media (
    id uuid NOT NULL ,
    post_id uuid NOT NULL,
    link TEXT NOT NULL,
    type TEXT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

