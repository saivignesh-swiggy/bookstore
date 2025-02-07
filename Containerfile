# Use Golang image to build the app
FROM golang:1.21 AS builder

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Enable CGO before building
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Build the Go application
RUN go build -o bookstore-app -ldflags '-extldflags "-static"'

# Use a minimal image to run the app
FROM alpine:latest

WORKDIR /root/

# Install sqlite dependencies (needed for CGO)
RUN apk add --no-cache sqlite-libs

# Copy compiled binary from the builder stage
COPY --from=builder /app/bookstore-app .

EXPOSE 8080

CMD ["./bookstore-app"]
