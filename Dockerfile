FROM golang:1.19-alpine as builder
RUN apk add --no-cache gcc libc-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main main.go

FROM alpine:latest as prod

WORKDIR /root
COPY --from=builder /app/main .
