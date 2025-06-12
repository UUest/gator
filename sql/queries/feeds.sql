-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
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
SELECT *
FROM feeds;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (feed_id, user_id)
VALUES (
$1,
$2
)
RETURNING *
)
SELECT
inserted_feed_follow.*,
feeds.name AS feed_name,
users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT
  ff.id,
  ff.created_at,
  ff.updated_at,
  ff.user_id,
  ff.feed_id,
  f.name AS feed_name,
  u.name AS user_name
FROM feed_follows ff
JOIN feeds f ON ff.feed_id = f.id
JOIN users u ON ff.user_id = u.id
WHERE ff.user_id = $1;

-- name: UnfollowFeed :one
DELETE FROM feed_follows
WHERE feed_id = $1 AND user_id = $2
RETURNING *;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT f.*
FROM feeds f
JOIN feed_follows ff ON f.id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY f.last_fetched ASC NULLS FIRST
LIMIT 1;
