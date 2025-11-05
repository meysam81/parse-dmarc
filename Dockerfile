# Build stage
FROM node:18-alpine AS frontend-builder

WORKDIR /build/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Go build stage
FROM golang:1.21-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /build/internal/api/dist ./internal/api/dist

RUN go build -o parse-dmarc ./cmd/parse-dmarc

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

COPY --from=backend-builder /build/parse-dmarc .

# Create directory for database
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/app/parse-dmarc"]
CMD ["-config=/data/config.json"]
