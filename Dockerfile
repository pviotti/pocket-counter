FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o pocket-counter .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/pocket-counter .

EXPOSE 8080

CMD ["./pocket-counter"]