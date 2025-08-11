-- name: CreatePost :one
INSERT INTO posts (
	id,
	created_at,
	updated_at,
	title,
	url,
	description,
	published_at,
	feed_id
) VALUES (
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
WITH feed_ids AS (
	SELECT feed_follows.feed_id AS id FROM feed_follows
	WHERE feed_follows.user_id = (
	SELECT users.id FROM users
	WHERE users.name = $1)
)  SELECT posts.* FROM posts
INNER JOIN
feed_ids
ON feed_ids.id = posts.feed_id;
