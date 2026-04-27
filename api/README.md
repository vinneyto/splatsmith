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

Standalone dev credentials (from `config/standalone.yaml`):

- username: `dev`
- password: `devpass`
- login endpoint: `POST /v1/auth/login`

## OpenAPI + code generation

- Source of truth: `api/internal/httpapi/openapi/openapi.yaml`
- Generator config: `api/internal/httpapi/openapi/oapi-codegen.yaml`
- Generated file: `api/internal/httpapi/openapi_gen.go`
- REST runtime: Gin router (`github.com/gin-gonic/gin`) with generated std-http handlers behind it.
- Public docs routes: `/docs` (redirect), `/scalar/*`, `/openapi.yaml`, `/openapi.json`
- Scalar API Reference (`scalarui/index.html`) and OpenAPI spec are embedded into the binary (`go:embed`).

Regenerate after spec changes:

```bash
cd api/internal/httpapi
go generate ./...
```

## Modes

- `standalone` — uses SQLite, idempotent async jobs, and simulated background execution for development
- `aws` — currently uses placeholders and returns not-implemented behavior for orchestration ports
