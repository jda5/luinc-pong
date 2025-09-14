
# ------------------------------------------------------------ build stage
# Use a Go image and give this stage a name: "builder"
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Only copy the dependency files first.
# Docker caches layers. As long as these two files don't change,
# Docker will re-use the cached layer from the next step, saving time.
COPY ./src/go.mod ./src/go.sum ./

# Download all dependencies into the container.
RUN go mod download

# Copy the rest of our source code.
COPY ./src .

# 'CGO_ENABLED=0' tells Go to build a "static binary," meaning
# it has ZERO external system dependencies. It's totally self-contained.
# '-o main' names the output file "main".
RUN CGO_ENABLED=0 go build -o main .

# ------------------------------------------------------------ run stage
# start fresh from a minimal OS image. Alpine is only ~5MB.
FROM alpine:latest

WORKDIR /app

# Copy the '/app/main' file from the "builder" stage
# into the current directory (.) of the new image.
COPY --from=builder /app/main .

# (Optional): Alpine is so minimal it doesn't
# even have the root CA certificates needed to make HTTPS calls.
# This line installs them, just in case we ever need them.
RUN apk add --no-cache ca-certificates

# Expose the port your Gin server is listening on
EXPOSE 8080

# Run the compiled program.
CMD ["./main"]