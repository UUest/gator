// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: posts.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :exec
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
)
`

type CreatePostParams struct {
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) error {
	_, err := q.db.ExecContext(ctx, createPost,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	return err
}

const getPostsForUser = `-- name: GetPostsForUser :many
SELECT p.id, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at, p.feed_id
FROM posts p
JOIN feed_follows ff ON p.feed_id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.updated_at DESC
LIMIT $2
`

type GetPostsForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
