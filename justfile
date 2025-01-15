set dotenv-load := true

up:
    docker run --rm \
        -v $(pwd)/cmd/migrate/migrations:/migrations \
        --network go-rag-poc_default \
        migrate/migrate \
        -path=/migrations/ \
        -database "postgresql://$RAG_DB_USERNAME:$RAG_DB_PASSWORD@$RAG_DB_HOST:$RAG_DB_PORT/$RAG_DB_DATABASE?sslmode=disable" \
        -verbose up

down COUNT='1':
    docker run --rm \
        -v $(pwd)/cmd/migrate/migrations:/migrations \
        --network go-rag-poc_default \
        migrate/migrate \
        -path=/migrations/ \
        -database "postgresql://$RAG_DB_USERNAME:$RAG_DB_PASSWORD@$RAG_DB_HOST:$RAG_DB_PORT/$RAG_DB_DATABASE?sslmode=disable" \
        -verbose down {{COUNT}}

to_version COUNT='1':
    docker run --rm \
        -v $(pwd)/cmd/migrate/migrations:/migrations \
        --network go-rag-poc_default \
        migrate/migrate \
        -path=/migrations/ \
        -database "postgresql://$RAG_DB_USERNAME:$RAG_DB_PASSWORD@$RAG_DB_HOST:$RAG_DB_PORT/$RAG_DB_DATABASE?sslmode=disable" \
        -verbose goto {{COUNT}}

force COUNT='1':
    docker run --rm \
        -v $(pwd)/cmd/migrate/migrations:/migrations \
        --network go-rag-poc_default \
        migrate/migrate \
        -path=/migrations/ \
        -database "postgresql://$RAG_DB_USERNAME:$RAG_DB_PASSWORD@$RAG_DB_HOST:$RAG_DB_PORT/$RAG_DB_DATABASE?sslmode=disable" \
        -verbose force {{COUNT}}

drop:
    docker run --rm \
        -v $(pwd)/cmd/migrate/migrations:/migrations \
        --network go-rag-poc_default \
        migrate/migrate \
        -path=/migrations/ \
        -database "postgresql://$RAG_DB_USERNAME:$RAG_DB_PASSWORD@$RAG_DB_HOST:$RAG_DB_PORT/$RAG_DB_DATABASE?sslmode=disable" \
        -verbose drop -f
