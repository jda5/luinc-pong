FROM golang:latest AS build

# Copy source code
COPY ./src ./app/src

WORKDIR /app/src

RUN go mod download
RUN go build -o main

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/src .

CMD ["./main"]