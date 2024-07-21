FROM golang:1.22-alpine

WORKDIR /runner

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /alertflow-runner

ENV RUNNER_ID=null
ENV ALERTFLOW_URL=null
ENV ALERTFLOW_APIKEY=null
ENV PAYLOADS_ENABLED=false
ENV PAYLOADS_PORT=8081
ENV PAYLOADS_MANAGERS=Alertmanager

EXPOSE ${PAYLOADS_PORT}

CMD [ "/alertflow-runner" ]
