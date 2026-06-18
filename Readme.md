# API Gateway

HTTP gateway for zenxbattle services. Routes gRPC backends via REST, compiles code via NATS.

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APIGATEWAYPORT` | No | `7000` | HTTP listen port |
| `ENVIRONMENT` | No | `development` | Runtime environment |
| `NATSURL` | No | `nats://localhost:4222` | NATS server URL |
| `USERGRPCURL` | No | `localhost:50051` | Auth gRPC endpoint |
| `PROBLEMGRPCURL` | No | `localhost:50055` | Problem gRPC endpoint |
| `FRONTENDURL` | No | `http://localhost:8080` | CORS allowed origin |
| `JWTSECRETKEY` | Yes | — | HMAC secret for JWT signing |
| `GOOGLECLIENTID` | No | — | Google OAuth client ID |
| `GOOGLECLIENTSECRET` | No | — | Google OAuth client secret |
| `GOOGLEREDIRECTURL` | No | — | Google OAuth callback |
