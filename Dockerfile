FROM golang:1.20-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o droneshield ./cmd

FROM alpine:latest

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/droneshield .

EXPOSE 8000

CMD ["./droneshield"]
