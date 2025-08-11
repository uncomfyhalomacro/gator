-- +goose Up
CREATE TABLE posts(
	id		UUID		PRIMARY KEY,
	created_at 	TIMESTAMP	NOT NULL,
	updated_at 	TIMESTAMP	NOT NULL,
	title		TEXT,
	url		TEXT,
	description	TEXT,
	published_at	TIMESTAMP,
	feed_id		UUID		NOT NULL
);

-- +goose Down
DROP TABLE posts;
