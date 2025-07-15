FROM golang:latest

# Build arguments for architecture-specific settings
ARG TARGETARCH
ARG TARGETPLATFORM
ARG DEBUG=false

# Set service name based on architecture
ENV SERVICE_NAME=service-pgn-partner-push-notification-linux-${TARGETARCH}

# Install Delve with architecture-specific settings
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest; \
    else \
        go install github.com/go-delve/delve/cmd/dlv@latest; \
    fi

RUN export CACHE_BUST=$(date +%s)

WORKDIR /app

# Copy the pre-compiled binaries
COPY app/${SERVICE_NAME} /app/

# Make binaries executable
RUN chmod +x /app/${SERVICE_NAME}

# Set environment variables
ENV TZ=Asia/Jakarta

# Use a shell script to determine whether to run in debug mode or not
CMD if [ "$DEBUG" = "true" ]; then \
        dlv --listen=0.0.0.0:40001 --headless=true --api-version=2 --accept-multiclient exec --continue /app/${SERVICE_NAME}; \
    else \
        /app/${SERVICE_NAME}; \
    fi 