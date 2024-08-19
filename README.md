# AlertFlow Runner

## Running Methods

### Master
All components are enabled. The runner will receive payloads, process them and scan for pending jobs.

### Worker
The Worker mode will disable the payload receiver component. The runner will only act as an Job executor.

### Listener
The runner will only act as a payload receiver. There will be no components enable to scan or execute any jobs.
