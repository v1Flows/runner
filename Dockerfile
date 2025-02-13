# Build tools stage
FROM golang:1.23 as build-tools

RUN mkdir -p /go-tools && \
    cp -r /usr/local/go /go-tools/ && \
    rm -rf /go-tools/go/pkg/*/cmd && \
    rm -rf /go-tools/go/pkg/bootstrap && \
    rm -rf /go-tools/go/pkg/obj && \
    rm -rf /go-tools/go/pkg/tool/*/api && \
    rm -rf /go-tools/go/pkg/tool/*/go_bootstrap

FROM golang:1.23 as builder
WORKDIR /app

# Install protobuf compiler and build dependencies
RUN apt-get update && \
    apt-get install -y \
    git \
    gcc \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Generate protobuf files
RUN protoc --go_out=pkg/plugin --go_opt=paths=source_relative \
    --go-grpc_out=pkg/plugin --go-grpc_opt=paths=source_relative \
    proto/plugin.proto

# Build the main application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o alertflow-runner ./cmd/alertflow-runner

# Final stage
FROM debian:bullseye-slim
WORKDIR /app

# Install runtime dependencies and Go tools
COPY --from=build-tools /go-tools/go /usr/local/go
ENV PATH="/usr/local/go/bin:${PATH}" \
    GOROOT="/usr/local/go" \
    GOPATH="/go"

# Install only git as runtime dependency
RUN apt-get update && \
    apt-get install -y \
    git \
    gcc \
    ca-certificates \
    protobuf-compiler && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p "$GOPATH/bin" && \
    chmod -R 755 "$GOPATH"

# Copy binary from builder
COPY --from=builder /app/alertflow-runner /app/alertflow-runner

# Create necessary directories
RUN mkdir -p /app/config /app/plugins && \
    chmod 755 /app/config /app/plugins

# Copy configuration
COPY config/config.yaml /app/config/config.yaml

VOLUME [ "/app/config", "/app/plugins" ]

EXPOSE 8081

# Set entrypoint
ENTRYPOINT ["/app/alertflow-runner"]
CMD ["-c", "/app/config/config.yaml"]