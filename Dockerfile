# FROM golang:latest

# WORKDIR /app

# COPY ./src /app/src

# WORKDIR /app/src

# RUN go mod download
# RUN go build -o main

# CMD ["./main"]

#-------------------------------------------------
# STAGE 1: The "Builder"
#-------------------------------------------------
# We use a Go image (based on lightweight Alpine Linux)
# and we give this stage a name: "builder"
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# This is a key optimization. We copy ONLY the dependency files first.
# Why? Docker caches layers. As long as these two files don't change,
# Docker will re-use the cached layer from the next step, saving time.
COPY ./src/go.mod ./src/go.sum ./

# Download all dependencies into the container.
RUN go mod download

# NOW we copy the rest of our source code.
# If we only change a .go file (and not go.mod), Docker will start
# from this cached step and skip the 'go mod download'.
COPY ./src .

# This is the build command.
# 'CGO_ENABLED=0' tells Go to build a "static binary," meaning
# it has ZERO external system dependencies. It's totally self-contained.
# '-o main' names the output file "main".
RUN CGO_ENABLED=0 go build -o main .

#-------------------------------------------------
# STAGE 2: The Final, Tiny Image
#-------------------------------------------------
# We start COMPLETELY fresh from a minimal OS image.
# Alpine is a popular choice because it's only ~5MB.
# Everything from STAGE 1 is now GONE...
FROM alpine:latest

WORKDIR /app

# ...EXCEPT for what we explicitly copy from that stage.
# This is the magic! We copy the '/app/main' file FROM the "builder" stage
# into the current directory (.) of our new, fresh image.
COPY --from=builder /app/main .

# (Optional but good practice): Alpine is so minimal it doesn't
# even have the root CA certificates needed to make HTTPS calls.
# This line installs them, just in case your app ever needs them.
RUN apk add --no-cache ca-certificates

# Expose the port your Gin server is listening on
EXPOSE 8080

# This is the command that will run when the container starts:
# just run our compiled program.
CMD ["./main"]