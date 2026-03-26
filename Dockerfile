ARG GO_IMG_VER=${GO_IMG_VER:-1.26.1-alpine}

FROM --platform=$BUILDPLATFORM golang:${GO_IMG_VER} AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

RUN adduser -D builduser
WORKDIR /app

COPY . .
RUN go mod download

## Build binaries

## Build API
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o /bin/serverbeacon-api ./cmd/api

## Build CLI
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o /bin/serverbeacon ./cmd/cli

## API RUNTIME
FROM gcr.io/distroless/static-debian12 AS serverbeacon-api

WORKDIR /app
COPY --from=builder /bin/serverbeacon-api /app/serverbeacon-api
USER 1000
EXPOSE 18080
ENTRYPOINT ["/app/serverbeacon-api"]

## CLI RUNTIME
FROM gcr.io/distroless/static-debian12 AS serverbeacon

WORKDIR /app
COPY --from=builder /bin/serverbeacon /app/serverbeacon
USER 1000
ENTRYPOINT ["/app/serverbeacon"]
