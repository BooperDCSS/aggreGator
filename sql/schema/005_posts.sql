-- +goose Up

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    published_at TEXT,
    title TEXT,
    url TEXT NOT NULL,
    description TEXT,
    feed_id UUID NOT NULL,
);


-- +goose Down

DROP TABLE posts;
