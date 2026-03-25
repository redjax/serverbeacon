ARG GO_IMG_VER=${GO_IMG_VER:-1.26.1-alpine}

FROM golang:${GO_IMG_VER} AS builder

RUN adduser -D builduser
WORKDIR /app

COPY . .
RUN go mod download

## Build binaries

## Build API
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/serverbeacon-api cmd/api/main.go

## Build CLI
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/serverbeacon cmd/cli/main.go

## API RUNTIME
FROM gcr.io/distroless/static-debian12 AS serverbeacon-api

WORKDIR /app
COPY --from=builder /bin/serverbeacon-api /app/
USER 1000
EXPOSE 18080
ENTRYPOINT ["./serverbeacon-api"]

## CLI RUNTIME
FROM gcr.io/distroless/static-debian12 AS serverbeacon

WORKDIR /app
COPY --from=builder /bin/serverbeacon /app/
USER 1000
ENTRYPOINT ["./serverbeacon", "start", "api"]
