# Runner

### What's the purpose?

The Runner is periodically checking the alertflow api for new Jobs to run.

### How does he report?

The Alertflow API exposes and API endpoint which catches and job progression statuses.

### Plugins?

The runner will open an HTTP Port which can be changed later via Env, Parameter or yaml config.
This Port will publish the following Endpoints:

| URL                        | Method | Description                        |
| -------------------------- | ------ | ---------------------------------- |
| /runner/status             | GET    | Get Status of the Runner           |
| /runner/publish/action     | POST   | Publish an Action to the Runner    |
| /runner/publish/payload    | POST   | Publish an new Payload             |
| /runner/executions/pending | GET    | Get Pending Execution on Alertflow |
| /runner/executions/finish  | POST   | Mark Execution as Finished         |
| /runner/executions/start   | POST   | Mark Execution as Started          |
