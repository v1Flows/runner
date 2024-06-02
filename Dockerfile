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
ENV PLUGIN_ENABLE=false
ENV PLUGIN_PORT=9854

EXPOSE ${PLUGIN_PORT}

CMD [ "/alertflow-runner" ]
