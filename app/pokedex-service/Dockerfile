FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/pokedex-service ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/pokedex-service /usr/local/bin/pokedex-service
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/pokedex-service"]
