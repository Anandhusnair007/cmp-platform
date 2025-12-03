# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service (service name passed as build arg)
ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/${SERVICE} ./cmd/${SERVICE}

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
ARG SERVICE
COPY --from=builder /build/bin/${SERVICE} /app/${SERVICE}

EXPOSE 8080

CMD ["/app/${SERVICE}"]
