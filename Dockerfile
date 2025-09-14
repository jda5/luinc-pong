FROM golang:latest

WORKDIR /app

COPY ./src /app/src

WORKDIR /app/src

RUN go mod download
RUN go build -o main

CMD ["./main"]