ARG GO_IMG_VER=${GO_IMG_VER:-1.26.1}
FROM golang:${GO_IMG_VER}

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /serverbeacon cmd/serverbeacon/main.go

EXPOSE 18080

CMD ["/serverbeacon"]
