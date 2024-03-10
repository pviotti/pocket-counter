FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --update --no-cache gcc musl-dev

COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY . .
RUN go build -o pocket-counter .

FROM alpine:3.20

WORKDIR /app

VOLUME /app/data

COPY static /app/static
COPY --from=builder /app/pocket-counter .

EXPOSE 8080

CMD ["./pocket-counter"]