# Build stage
FROM golang:1.25-alpine@sha256:1ae0735f00daffa3aaf1363a5184c0d2dc55c78e3db4ec70241cdac97bf84b59 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mobile-bff ./cmd/server

# Runtime stage
FROM alpine:3.20

RUN adduser -D -u 1000 appuser

WORKDIR /home/appuser/

COPY --from=builder /app/mobile-bff .

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

USER appuser

CMD ["./mobile-bff"]
