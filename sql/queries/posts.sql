-- name: CreatePost :one

INSERT INTO posts (id, created_at, updated_at, published_at, title, url, description, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many

SELECT
    posts.tile AS post_title,
    posts.url AS post_url
FROM posts
    INNER JOIN feeds_follows ON posts.feed_id = feeds_follows.feed_id
WHERE
    feeds_follows.user_id = $1
LIMIT
    $2
ORDER BY
    posts.published_at DESC;
