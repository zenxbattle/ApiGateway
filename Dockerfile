FROM golang:1.24.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git
RUN --mount=type=secret,id=GIT_AUTH_TOKEN \
    git config --global url."https://oauth2:$(cat /run/secrets/GIT_AUTH_TOKEN)@github.com/".insteadOf "https://github.com/"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o /app/main ./cmd

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
