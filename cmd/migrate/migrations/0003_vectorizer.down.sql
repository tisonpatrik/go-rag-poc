BEGIN;

SELECT ai.drop_vectorizer(1, drop_all=>true);

COMMIT;