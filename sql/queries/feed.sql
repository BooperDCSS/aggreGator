-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many

SELECT * FROM feeds;

-- name: IDFromURL :one

SELECT id FROM feeds
WHERE url = $1;

-- name: CreateFeedFollow :one

WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
FROM inserted_feed_follow
    INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
    INNER JOIN users ON inserted_feed_follow.user_id = users.id
ORDER BY
    users.name;

-- name: GetFeedFollowsForUser :many

SELECT
    feeds.name AS feed_name
FROM feed_follows
    INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
    feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec

DELETE FROM feed_follows
WHERE
    feed_follows.user_id = $1
    AND
    feed_follows.feed_id = $2;

-- name: MarkFeedFetched :exec

UPDATE
    feed_follows
SET
    updated_at = transaction_timestamp(),
    last_fetched_at = transaction_timestamp()
WHERE
    feed_id = $1;

-- name: GetNextFeedToFetch :one

SELECT
    url
FROM feeds
    INNER JOIN feed_follows ON feeds.id = feed_follows.feed_id
ORDER BY
    feed_follows.last_fetched_at NULLS FIRST;
