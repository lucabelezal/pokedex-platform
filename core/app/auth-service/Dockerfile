FROM golang:1.25-alpine@sha256:e39446288ca670e24c72e6e285574e314d10ad0978c3b033b6273c83e9cbeed3 AS builder
WORKDIR /src

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/auth-service ./cmd/server

FROM alpine:3.20
RUN adduser -D -u 1000 appuser
WORKDIR /app
COPY --from=builder /bin/auth-service /usr/local/bin/auth-service
EXPOSE 8082
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1
USER appuser
ENTRYPOINT ["/usr/local/bin/auth-service"]
