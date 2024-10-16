FROM golang:1.22-alpine as builder 

WORKDIR /runner

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /alertflow-runner

FROM alpine:3.12 as runner

COPY --from=builder /alertflow-runner /alertflow-runner

RUN mkdir /runner
RUN mkdir /runner/config
COPY handlers/config/config.yaml /runner/config/config.yaml

VOLUME [ "/runner" ]

EXPOSE 8081

CMD [ "/alertflow-runner", "-c", "/runner/config/config.yaml" ]
