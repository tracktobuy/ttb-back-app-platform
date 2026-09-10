# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

`ttb-back-app-platform` (module `github.com/tracktobuy/ttb-back-app-platform`) is the Go backend API for the
TrackToBuy app. It's a REST API backed by MongoDB with no external web framework — routing uses the stdlib
`http.ServeMux` (Go 1.22+ method+pattern routes), not the `httprouter` dependency in `go.mod` (that package is
only referenced by a vestigial, unused `helper.ReadParam` function).

## Commands

```bash
go run ./cmd/api              # run the API locally (reads .env via godotenv, falls back to real env vars)
go build -o server ./cmd/api  # build the binary
go vet ./...                  # static analysis
go test ./...                 # run tests (no test files currently exist in the repo)
gofmt -l .                    # check formatting
```

There is no Makefile, linter config, or CI test step — `.github/workflows/1_deploy.yml` builds and deploys the
Docker image straight from `main` with no test gate.

Local setup: copy `.env.example` to `.env` and fill in `MONGO_USERNAME`/`MONGO_PASSWORD`/`MONGO_URI`/`MONGO_DB` (a
running MongoDB Atlas-style cluster is required — `config.LoadConfig` calls `log.Fatal` if Mongo credentials are
missing, and `config.MongoConnect` fatals if it can't ping the cluster). `.vscode/launch.json` has a "Debug TTB API"
configuration that runs `cmd/api` with `.env` loaded.

Docker:
```bash
docker build -t ttb-backend-api .
docker run --name ttb-back-service -e MONGO_HOST=... -e MONGO_DB_USER=... -e MONGO_DB_PASSWORD=... \
  -e MONGO_DB_NAME=... -e API_SERVER_PORT=9999 -p 8080:9999 ttb-backend-api
```

### Environment variables

| Variable | Description | Required |
| --- | --- | --- |
| `MONGO_HOST` | Hostname (and connection options) of the MongoDB cluster | Yes |
| `MONGO_DB_USER` | Username used to authenticate with the MongoDB database | Yes |
| `MONGO_DB_PASSWORD` | Password used to authenticate with the MongoDB database | Yes |
| `MONGO_DB_NAME` | Name of the MongoDB database to connect to | Yes |
| `API_SERVER_PORT` | Port on which the API server listens. Default: 8080 | No |
| `CORS_ALLOWED_ORIGINS` | Comma separated list of allowed origins | No |

Note the mismatch between `.env.example` (`MONGO_USERNAME`/`MONGO_PASSWORD`/`MONGO_URI`) and what `config.go`
actually reads (`MONGO_DB_USER`/`MONGO_DB_PASSWORD`/`MONGO_HOST`) — trust `config/config.go` as source of truth.

## Architecture

Strict layering: **handler → service → repository → MongoDB**, with each layer defined by a Go interface and a
private struct implementation, and dependencies wired top-down through two composition roots:

- `internal/repository.go` (`internal.CreateRepositories`) builds every `repository.XxxRepository` from a
  `*mongo.Database`.
- `internal/service.go` (`internal.CreateServices`) builds every `service.XxxService` from the `Repository` struct.
  Some services compose other services rather than repositories directly (see `AccountService` below).
- `cmd/api/handler.go` (`CreateHandlers`) builds every `handler.XxxHandler` from the full `internal.Service` struct
  (handlers get the whole service bundle, not just the one they're named after — e.g. `ItemHandler.Create` also
  calls `GroupService`, `StoreService`, and `GroupItemService` to build a new item, its price/store record, and the
  group-item junction row in one request).

`cmd/api/main.go` wires it all together: loads config → connects to Mongo → builds handlers → starts the server
(`cmd/api/server.go`). Routes are registered in `cmd/api/routes.go` using Go's `"METHOD /path/{param}"` mux syntax,
wrapped in the CORS middleware from `cmd/api/middleware.go` (`enableCORS` — origin allow-list from
`CorsAllowedOrigins`, handles `OPTIONS` preflight itself since the mux only registers concrete methods).

### Package layout (`internal/`)

- `domain/` — MongoDB document structs (`bson`/`json` tags), one type per collection (`User`, `Group`, `Item`,
  `Store`, `UserGroup`, `GroupItem`). Every domain type carries `UUID` (public-facing, UUIDv7 via
  `helper.GenerateUUIDV7`) and Mongo `_id` `primitive.ObjectID` (internal joins/lookups) side by side, plus a
  `Version` field (present but not currently enforced as optimistic-concurrency control anywhere).
- `dto/request/`, `dto/response/` — API-facing shapes, kept separate from `domain` types; services translate
  between them (e.g. `itemService.formatItemResponse`). `dto/cookie/` defines the `Account` struct stored in the
  auth cookie.
- `repository/` — one file per aggregate, plus `colletions.go` (name has a typo, keep it when referencing) which
  centralizes Mongo collection name constants. Repositories hand-build `bson.M` filters and `mongo.Pipeline`
  aggregations directly (see `item_repository.go`'s `GetAllByUserId`, which joins `user_group → users → groups →
  group_item → items` to resolve a user's items).
- `service/` — business logic and request/response shaping. `crud_service.go` defines a generic
  `CrudService[T any]` interface, but individual services largely implement their own bespoke method sets rather
  than that generic interface. `AccountService` is the one service that orchestrates other services instead of a
  repository: `CreateAccount` creates a `User`, then a default `Group` for them, then the `UserGroup` join row.
- `handler/` — HTTP handlers implementing `net/http` signatures directly (no framework request/response types).
  `envelope.go` defines the shared `envelope map[string]any` used to shape JSON bodies as `{"data": ..., "error":
  ..., "message": ...}`.
- `helper/` — shared HTTP utilities: `WriteJSON`/`ReadJSON`, error responders (`BadRequest`, `NotFound`,
  `InternalServerError`) that build `response.ClientError`, `DateTime` (formats timestamps in `America/Sao_Paulo`),
  UUIDv7 generation, and the auth cookie helpers below.
- `logger/` — a thin `log/slog` wrapper (`logger.Logger`). Call `NewLogger()` then tag it with
  `SetHandlerName`/`SetServiceName`/`SetRepositoryName`/`SetMiddlewareName`/`SetConfigName`/`SetComponentName` and
  `SetMethodName` before logging, so every JSON log line carries `handler`/`service`/`repository`/`method` fields
  for tracing which layer and function produced it. Follow this pattern (tag with layer + component name in the
  constructor, `SetMethodName` at the top of each method) when adding new handlers/services/repositories.

### Auth model

There's no token-based auth. `helper.SetCookie`/`helper.GetCookie` store/read a base64-encoded JSON blob (the
`dto/cookie.Account` struct) in an HTTP-only `account` cookie. Handlers call `helper.GetCookie(w, r)` to identify
the current user, then look up the full `domain.User` via `UserService.GetById`. There's no `POST /login` or
`/sessions` route yet — nothing currently calls `helper.SetCookie`; account creation (`POST /accounts` →
`AccountHandler.CreateAccount` → `AccountService.CreateAccount`) still needs to be wired to actually issue the
cookie for a new session.

### Deployment

`Dockerfile` is a two-stage build: `golang:1.26-alpine` compiles a static binary (`CGO_ENABLED=0`), then it's
copied into `gcr.io/distroless/static:nonroot`. `.github/workflows/1_deploy.yml` runs on push to `main`: builds and
pushes the image to AWS ECR (OIDC auth, no long-lived keys), then SSHes into a VPS to `docker compose pull && up
-d` using env values written into a generated `.env`. The `k8s/` manifests (`configmap.yaml`, `deployment.yaml`,
`secrets.yaml`) exist but are not what the current deploy workflow uses.
