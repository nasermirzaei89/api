-- +migrate Up

CREATE TABLE users
(
    uuid          TEXT NOT NULL PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

-- +migrate Down

DROP TABLE users CASCADE;
