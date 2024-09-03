# Multi-stage go image build

# STAGE 1: Build the executable
FROM golang:1.22.5 AS builder

# Set the working directory
WORKDIR /build

# Copy the source code
COPY . .

# Download and install dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 go build -o /build/kerigma

# STAGE 2: Create the final minimal image
FROM scratch

# Copy the built executable from the builder stage
COPY --from=builder /build/kerigma /kerigma

# Copy CA certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Expose the app port
EXPOSE 8080

# Run the executable
CMD ["/kerigma"]