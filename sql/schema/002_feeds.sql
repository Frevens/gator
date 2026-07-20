-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

ALTER TABLE feeds
ADD CONSTRAINT fk_feeds_users
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE;

-- +goose Down
ALTER TABLE feeds DROP CONSTRAINT IF EXISTS fk_feeds_users;
DROP TABLE IF EXISTS feeds;