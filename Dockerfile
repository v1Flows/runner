FROM golang:1.22-alpine as builder 

WORKDIR /runner

# Update the package list and install git and gcc
RUN apk update && apk add --no-cache git gcc

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /alertflow-runner ./cmd/alertflow-runner

RUN mkdir -p /runner/config
COPY config/config.yaml /runner/config/config.yaml

VOLUME [ "/runner" ]

EXPOSE 8081

CMD [ "/alertflow-runner", "-c", "/runner/config/config.yaml" ]
