FROM golang:1.23 as builder

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

VOLUME [ "/app/config", "/app/plugins" ]

EXPOSE 8081

CMD [ "/runner", "-c", "/app/config/config.yaml" ]