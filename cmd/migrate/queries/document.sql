-- name: InsertDocument :one
INSERT INTO document (
    id, 
    document_name, 
    date_time, 
    original_link, 
    content
)
VALUES (
    uuid_generate_v4(),
    $1,
    $2,
    $3,
    $4)
RETURNING *;
