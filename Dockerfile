ARG GO_IMG_VER=${GO_IMG_VER:-1.26.1-alpine}

FROM golang:${GO_IMG_VER} AS builder

RUN adduser -D builduser

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o serverbeacon cmd/serverbeacon/main.go

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app .

USER 1000

EXPOSE 18080

CMD ["./serverbeacon"]
