-- +goose Up
CREATE TABLE feed_follows (
    id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT feed_follows_pk PRIMARY KEY (id),
    CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT feeds_fk FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
    CONSTRAINT unique_feed_user UNIQUE (feed_id, user_id)
);

-- +goose Down
DROP TABLE feed_follows;
