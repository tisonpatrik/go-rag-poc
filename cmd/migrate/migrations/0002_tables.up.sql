BEGIN;

CREATE TABLE document (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_name VARCHAR NOT NULL,
    date_time TIMESTAMPTZ NOT NULL,
    original_link TEXT NOT NULL,
    content TEXT NOT NULL
);

COMMIT;
