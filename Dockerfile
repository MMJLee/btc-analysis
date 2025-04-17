FROM golang:1.24.1-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/candle


FROM alpine:latest
COPY --from=builder /bin/candle /bin/candle
CMD ["/bin/candle --ticker BTC-USD --mode track"]