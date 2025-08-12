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
) 
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

-- name: GetPostByUser :many
SELECT
    p.*,
    f.name as feed_name
FROM 
    posts AS p
INNER JOIN 
    feeds AS f
    ON p.feed_id = f.id
INNER JOIN 
    users AS u
    ON f.user_id = u.id
WHERE 
    u.name = $1
ORDER BY 
    p.created_at DESC
LIMIT $2
;
