FROM golang:1.25-alpine@sha256:1ae0735f00daffa3aaf1363a5184c0d2dc55c78e3db4ec70241cdac97bf84b59 AS builder
WORKDIR /src

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/pokemon-catalog-service ./cmd/server

FROM alpine:3.20
RUN adduser -D -u 1000 appuser
WORKDIR /app
COPY --from=builder /bin/pokemon-catalog-service /usr/local/bin/pokemon-catalog-service
EXPOSE 8081
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1
USER appuser
ENTRYPOINT ["/usr/local/bin/pokemon-catalog-service"]
