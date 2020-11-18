-- +migrate Up

CREATE TABLE posts
(
    uuid             TEXT NOT NULL PRIMARY KEY,
    title            TEXT NOT NULL,
    slug             TEXT NOT NULL UNIQUE,
    content_markdown TEXT NOT NULL,
    content_html     TEXT NOT NULL
);

-- +migrate Down

DROP TABLE posts CASCADE;
