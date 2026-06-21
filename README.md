# ZenXBattle API Gateway

REST-to-gRPC gateway for the ZenXBattle platform. Single entrypoint for the frontend — handles auth, routing, rate limiting, and CORS.

## Architecture

```
Browser → ApiGateway (REST, :8080) → gRPC → AuthUserService
                                    → gRPC → ProblemService
                                    → gRPC → ChallengeService
                                    → gRPC → CodeExecutionEngine
```

## Tech Stack

- **Go** + Gin (HTTP framework)
- **gRPC** clients to all downstream services
- **Prometheus** metrics (`/metrics`)
- **Ristretto** in-memory cache
- **Zap** structured logging
- **Rate limiting** per IP/endpoint

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/metrics` | Prometheus metrics |
| `POST` | `/api/auth/*` | Auth routes (proxied to AuthUserService) |
| `POST` | `/api/problems/*` | Problem CRUD |
| `GET` | `/api/leaderboard/*` | Leaderboard queries |
| `POST` | `/api/challenges/*` | Challenge creation/join/submit |
| `GET` | `/api/execute/*` | Code execution |

## Quick Start

```bash
# Set environment
export AUTH_SERVICE_ADDR=localhost:50051
export PROBLEM_SERVICE_ADDR=localhost:50052
export CHALLENGE_SERVICE_ADDR=localhost:50053
export CODE_ENGINE_ADDR=localhost:50054

# Run
go run cmd/main.go
# → Gateway listening on :8080
```

## Config

See `configs/` — uses environment variables with sensible defaults.

## Related Services

- [AuthUserAdminService](https://github.com/zenxbattle/AuthUserAdminService) — user auth & admin gRPC
- [ProblemService](https://github.com/zenxbattle/ProblemService) — problem CRUD gRPC
- [ChallengeService](https://github.com/zenxbattle/ChallengeService) — real-time battle WebSocket + gRPC
- [CodeExecutionEngine](https://github.com/zenxbattle/CodeExecutionEngine) — sandboxed code execution
- [Frontend](https://github.com/zenxbattle/Frontend) — React SPA
- [CommonProto](https://github.com/zenxbattle/CommonProto) — shared protobuf definitions

## Docker

```bash
docker build -t zenxbattle-api-gateway .
docker run -p 8080:8080 zenxbattle-api-gateway
```

## Deploy (K3s)

```bash
kubectl apply -k https://github.com/zenxbattle/infrastructure/tree/k3s/k3s/services/api-gateway
```
