FROM golang:1.27.3-alpine3.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/ ./cmd/
COPY ./pkg/ ./pkg/
COPY ./internal/ ./internal/
COPY ./configs/container.config.yaml ./configs/config.yaml

RUN CGO_ENABLED=0 GOOS=linux go build -o auth_bin ./cmd/app/main.go

FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/auth_bin .
COPY --from=builder /app/configs/config.yaml /app/configs/config.yaml

CMD [ "./auth_bin" ]