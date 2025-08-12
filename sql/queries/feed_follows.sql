-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    ff.*,
    f.name AS feed_name,
    u.name AS user_name
FROM 
    inserted_feed_follow AS ff
INNER JOIN 
    users AS u
    ON ff.user_id = u.id
INNER JOIN 
    feeds AS f
    ON ff.feed_id = f.id
;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.*,
    f.name AS feed_name
FROM 
    feed_follows AS ff
INNER JOIN 
    feeds AS f
    ON ff.feed_id = f.id
WHERE 
    ff.user_id = $1
;

-- name: DeleteFeedFollow :exec
DELETE
FROM 
    feed_follows AS ff
USING users AS u, feeds AS f
WHERE 
    ff.feed_id = f.id
    AND ff.user_id = u.id
    AND u.name = $1
    AND f.url = $2
;
