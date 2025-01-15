-- name: InsertDocument :one
INSERT INTO document (
    id, 
    document_name, 
    date_time, 
    original_link, 
    html_content,
    content,
    doc_language
)
VALUES (
    uuid_generate_v4(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;
