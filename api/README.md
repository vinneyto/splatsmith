# api

Go backend skeleton split into modules:

- `internal/core` — business-domain contracts and minimal auth logic
- `internal/standalone` — local implementations (SQLite + local filesystem + stubs)
- `internal/aws` — AWS adapter stubs (to be implemented later)
- `cmd/api` — binary entrypoint with mode-based bootstrap

## Run

```bash
cd api
go run ./cmd/api -config ./config/standalone.yaml
```

## Modes

- `standalone` — uses SQLite and local adapters for development
- `aws` — currently uses placeholders and returns not-implemented behavior
