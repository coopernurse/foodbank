# Use an official lightweight Go image
FROM golang:1.23-alpine AS builder

# Set up working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Start fresh with a new image
FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /root/

# Copy the built executable
COPY --from=builder /app/main .
COPY --from=builder /app/static/ static/

# Expose port and set entrypoint
EXPOSE 8080
ENTRYPOINT ["./main"]
