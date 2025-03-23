FROM golang:1.24 as builder

WORKDIR /app

# Update the package list and install git and gcc
RUN apt-get update && apt-get install -y git gcc

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build the main application
RUN CGO_ENABLED=1 GOOS=linux go build -o /runner ./cmd/runner

# Copy configuration files
RUN mkdir -p /app/config
RUN mkdir -p /app/plugins
COPY config/config.yaml /app/config/config.yaml

# Install Ansible in the final stage
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    python3 python3-pip sshpass && \
    pip3 install --no-cache-dir ansible --break-system-packages && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

VOLUME [ "/app/config", "/app/plugins" ]

EXPOSE 8081

CMD [ "/runner", "-c", "/app/config/config.yaml" ]