BEGIN;

-- Configure the vectorizer to vectorize `content` with metadata formatted for embedding
SELECT ai.create_vectorizer(
   'document'::regclass,    -- Source table
   embedding => ai.embedding_openai(
       'text-embedding-3-small', -- Embedding model
       1536                      -- Embedding dimensions
   ),
   chunking => ai.chunking_recursive_character_text_splitter(
       'content',                -- Column to chunk
       chunk_size => 1200,       -- Max characters per chunk
       chunk_overlap => 150,     -- Overlap for context
       separators => array[E'\n\n', E'\n', '.', ' '] -- Recursive splitting strategy
   ),
   formatting => ai.formatting_python_template(
        'Sequence Document Name: $document_name\n' ||
        'Content: $chunk'
   )
);

COMMIT;
