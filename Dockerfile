# Start from the official Go base image
FROM golang:1.22-alpine as builder

# Install make and other necessary build tools
RUN apk add --no-cache make gcc musl-dev

# Copy files for building application
WORKDIR /app
COPY go.mod go.sum ./
COPY . .

# Install, generate docs and build application
RUN make install
RUN make docs
RUN make build

# Use a minimal base image for the final build
FROM alpine:latest

# Install necessary CA certificates
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/bin/liteapi .

# Expose the port the application runs on
EXPOSE $PORT

CMD ["./liteapi"]
