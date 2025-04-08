# Use the official Golang image as the base image
FROM golang:alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install just
RUN apk add --no-cache just upx

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN just build

# Compress the binary with UPX
RUN upx --lzma main

# Start a new stage from scratch
FROM scratch AS runner

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

ENTRYPOINT ["/app/main"]