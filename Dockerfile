FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o riffle ./cmd/riffle

# Final stage
FROM alpine:latest

WORKDIR /app
# Create configuration directory
RUN mkdir -p /app/conf

COPY --from=builder /app/riffle /app/riffle

ENTRYPOINT ["/app/riffle"] 