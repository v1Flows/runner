FROM golang:1.22-alpine

WORKDIR /runner

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /alertflow-runner

ENV LOG_LEVEL=Info
ENV RUNNER_ID=null
ENV MODE=master
ENV ALERTFLOW_URL=null
ENV ALERTFLOW_API_KEY=null
ENV PAYLOADS_ENABLED=true
ENV PAYLOADS_PORT=8080

EXPOSE ${PLUGIN_PORT}

CMD [ "/alertflow-runner", "--config.file=config.yaml" ]
