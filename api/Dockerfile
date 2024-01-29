# syntax=docker/dockerfile:1

FROM golang:1.20.2-alpine as builder

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

COPY ./ .

RUN go mod download

RUN go build -o /app/start.sh ./cmd/serverd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/start.sh ./

RUN chmod +x start.sh

CMD ["./start.sh"]