FROM golang:1.22

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY . .
RUN go build -o pocket-counter .

VOLUME /app/data

EXPOSE 8080

CMD ["./pocket-counter"]