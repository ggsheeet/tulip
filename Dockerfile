# Multi-stage go image build

# STAGE 1: Build the executable
FROM golang:1.23 AS builder
WORKDIR /build
COPY . .
RUN go mod download
# RUN CGO_ENABLED=0 go build -o /build/kerigma
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /build/kerigma

# STAGE 2: Development stage (for live reload)
FROM golang:1.23 AS dev
WORKDIR /app
COPY --from=builder /build .
RUN go install github.com/air-verse/air@latest
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# STAGE 3: Create the final minimal image for production
FROM scratch AS production
COPY --from=builder /build/kerigma /kerigma
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# EXPOSE 8080
CMD ["/kerigma"]
