-- +migrate Up

ALTER TABLE posts
    ADD COLUMN published_at TIMESTAMPTZ NULL;

-- +migrate Down

ALTER TABLE posts
    DROP COLUMN published_at;
