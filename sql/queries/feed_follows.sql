-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
	INSERT INTO feed_follows (
		id,
		created_at,
		updated_at,
		user_id,
		feed_id
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)
	RETURNING *
) SELECT inserted_feed_follow.*,
	feeds.name AS feed_name,
	users.name AS user_name
FROM inserted_feed_follow
INNER JOIN 
feeds 
ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN 
users 
ON users.id = inserted_feed_follow.user_id;

-- name: GetAllFeedFollows :many
SELECT * FROM feed_follows;

-- name: GetFeedFollowsForUser :many
WITH feed_ids AS (
	SELECT feed_follows.feed_id AS id FROM feed_follows
	WHERE feed_follows.user_id = (
	SELECT users.id FROM users
	WHERE users.name = $1)
)  SELECT feeds.name AS feed_name, feeds.url AS feed_url FROM feeds
INNER JOIN
feed_ids
ON feed_ids.id = feeds.id;

-- name: UnfollowFeedForUser :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = (
	SELECT users.id FROM users
	WHERE users.name = $1 LIMIT 1
) AND feed_follows.feed_id = (
	SELECT feeds.id FROM feeds
	WHERE feeds.name = $2
);
