FROM golang:1.19-alpine
RUN apk add --no-cache make gcc libc-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o out/main.out main.go
