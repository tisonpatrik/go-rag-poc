# go-rag-poc

Small PoC of RAG techniques written in Go

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Makefile Commands

### Run build make command with tests
```bash
make all
```

### Build the application
```bash
make build
```

### Run the application
```bash
make run
```

### Create DB container
```bash
make docker-run
```

### Shutdown DB Container
```bash
make docker-down
```

### DB Integrations Test
```bash
make itest
```

### Live reload the application
```bash
make watch
```

### Run the test suite
```bash
make test
```

### Clean up binary from the last build
```bash
make clean
```

## Justfile Commands

### Apply all migrations
```bash
just up
```

### Rollback the last migration
```bash
just down COUNT=1
```
Replace `COUNT` with the number of migrations to rollback.

### Migrate to a specific version
```bash
just to_version COUNT=1
```
Replace `COUNT` with the desired migration version.

### Force a specific migration version
```bash
just force COUNT=1
```
Replace `COUNT` with the migration version to force.

### Drop all migrations
```bash
just drop
```

### Note
Make sure to set the environment variables (`RAG_DB_USERNAME`, `RAG_DB_PASSWORD`, `RAG_DB_HOST`, `RAG_DB_PORT`, `RAG_DB_DATABASE`) properly before running the `Justfile` commands.
