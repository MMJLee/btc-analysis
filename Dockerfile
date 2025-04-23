FROM golang:1.24.1-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app

FROM alpine:latest

RUN apk add --no-cache redis
COPY --from=builder /bin/app /bin/app

CMD sh -c "redis-server & sleep 1 && exec /bin/app"