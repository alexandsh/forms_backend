FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

ENV GOPROXY=off

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o server ./cmd/server

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /app/server

ENV PORT=8080
EXPOSE 8080

CMD ["/app/server"]

