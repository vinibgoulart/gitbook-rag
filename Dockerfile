### Build
FROM golang:1.22.2 as builder

WORKDIR /app-todo-list

ENV GOOS=linux \
    GOARCH=amd64 

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build ./cmd/api

### Certs
FROM alpine:latest as certs
RUN apk --update add ca-certificates
RUN apk add --no-cache tzdata

### App
FROM ubuntu:bionic as app
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder app-todo-list/api /api

EXPOSE 8054

ENTRYPOINT [ "/api" ]
