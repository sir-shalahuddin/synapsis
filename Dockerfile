FROM golang:1.21.12-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN go build -o syn ./cmd/

FROM alpine 

WORKDIR /app

COPY --from=builder /app/syn ./syn

CMD ["./syn"]
