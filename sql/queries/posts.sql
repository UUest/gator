-- name: CreatePost :exec
INSERT INTO posts (
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
) VALUES (
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: GetPostsForUser :many
SELECT p.*
FROM posts p
JOIN feed_follows ff ON p.feed_id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.updated_at DESC
LIMIT $2;
