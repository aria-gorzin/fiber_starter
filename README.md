
# PIM – Product Information Management Service

## Table of contents
1. [Prerequisites](#prerequisites)
2. [Quick start (macOS)](#quick-start-macos)
3. [Infrastructure](#infrastructure)
4. [Database migrations](#database-migrations)
5. [Code generation](#code-generation)
6. [Running & testing](#running--testing)
7. [Proto / gRPC tooling](#proto--grpc-tooling)
8. [Utilities](#utilities)
9. [Kubernetes deployment](#kubernetes-deployment)
10. [Environment variables](#environment-variables)

---

## Prerequisites
| Tool | Purpose | Install |
|------|---------|---------|
| **Docker Desktop** | Containers & Compose | <https://www.docker.com/products/docker-desktop> |
| **Homebrew** | Package manager | <https://brew.sh/> |
| **Golang ≥ 1.22** | Backend code | `brew install go` |
| **Node / npm** | DBML & dbdocs CLI | `brew install node` |
| **TablePlus** | Database client (optional) | <https://tableplus.com/> |

> **macOS users:** skip the manual steps and run `make install_mac` to pull every dependency (brew, npm & go installs).

---

## Quick start (macOS)

```bash
# 1. Install all local tooling
make install_mac          # one‑shot install

# 2. Spin up Postgres & a Docker network
make network              # one‑time
make postgres             # background container

# 3. Create the database & run migrations
make createdb
make migrateup            # or `make migrateup1`

# 4. Launch the server localy
docker compose -f "docker-compose-local.yaml" up -d      # starts any extra services in docker-compose.yml
make server               # swag init + go run main.go

# 5. Run server
docker compose up -d
```

---

## Infrastructure

| Target | Description |
|--------|-------------|
| `make network` | Create Docker network `pim-network` |
| `make postgres` | Run Postgres 14-alpine on port 5432 |
| `make mysql` | Run MySQL 8 container on port 3306 |
| `make redis` | Launch Redis 7 on port 6379 |

---

## Database migrations

| Target | Action |
|--------|--------|
| `make migrateup` / `migrateup1` | Apply all / one migration |
| `make migratedown` / `migratedown1` | Roll back all / one migration |
| `make new_migration name=<snake_case>` | Create timestamped up/down SQL pair |
| `make test_migration` | Drop → create → up → down (sanity check) |

---

## Code generation

| Step | Command | Output |
|------|---------|--------|
| DB schema (SQL) | `make db_schema` | `doc/schema.sql` |
| DB docs | `make db_docs password=<pw>` | Publishes dbdocs.io project |
| CRUD queries | `make sqlc` | Go code in `db/sqlc` |
| Go mocks | `make mock` | `db/mock/*.go`, `worker/mock/*.go` |

---

## Running & testing

| Target | Action |
|--------|--------|
| `make server` | Regenerate Swagger (`swag init`) and run `main.go` |
| `make test` | `go test -v -cover -short ./…` |

---

## Proto / gRPC tooling

| Target | What it does |
|--------|--------------|
| `make proto` | Compile `.proto` files into Go, gRPC‑Gateway & OpenAPI, then embed swagger via `statik` |
| `make evans` | Launch [evans](https://github.com/ktr0731/evans) REPL on `localhost:9090` |

---

## Utilities

| Tool | Install note |
|------|--------------|
| **Gomock** (`mockgen`) | Installed by `make install_mac` |
| **Migrate CLI** | Installed by `make install_mac` |
| **dbdocs / dbml2sql** | Installed by `make install_mac` |

---

## Kubernetes deployment

```bash
# NGINX Ingress controller
kubectl apply -f   https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.48.1/deploy/static/provider/aws/deploy.yaml

# cert-manager (v1.4.0)
kubectl apply -f   https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml
```

---

## Environment variables

| Var | Default | Notes |
|-----|---------|-------|
| `APP_NAME` | `pim` | Defined in `Makefile` |
| `DB_URL` | `postgresql://root:secret@localhost:5432/pim?sslmode=disable` | Override if you change credentials or ports |

---

### Happy shipping 🚀
