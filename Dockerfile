# Step 1: Build the Go binary
FROM golang:1.24.3-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gochan main.go

# Step 2: Minimal runtime image with ImageMagick
FROM alpine:latest

# Install certificates and ImageMagick
RUN  apk update && apk add --no-cache \
    ca-certificates \
    imagemagick \
    libwebp-tools \
    ttf-freefont \
    libjpeg-turbo \
    ffmpeg \
    libpng 

# Create a non-root user and group
RUN addgroup -S gochan && adduser -S gochan -G gochan

# Create app directory and set perms
RUN mkdir -p /app/static/content \
    && chown -R gochan:gochan /app

# Set working directory
WORKDIR /app

# Copy Go binary from builder
COPY --from=builder /app/gochan .

#Copy static assets (css/js)
COPY static/js /app/static/js
COPY static/stylesheets /app/static/stylesheets
COPY static/fontawesome /app/static/fontawesome

#Copy migrations
COPY migrations /app/migrations

#Copu config
COPY config/config.toml /app/config/config.toml

#Copy entrypoint script
COPY docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port
EXPOSE 8080

# Switch to non-root user
USER gochan

# Run entrypoint script and then the binary
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./gochan"]
