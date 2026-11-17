-- +goose Up

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    published_at TIMESTAMP,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    UNIQUE(url),
    description TEXT,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_postsfeedid FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE
);


-- +goose Down

DROP TABLE posts;
