# Runner

## Incoming Payloads Enabled:
1. Payload reached Runner
2. Post Payload on AlertFlow API and mark as already took by runner
3. Pull all Flows which are are created for that Project the runner is dedicated for
4. Check Action by Action for pattern agains the payload
?. what is multiple Action match?
5. Execute Action where Pattern matches
6. Post Updates about that running Action on AlertFlow API
7. Post Finish Update on AlertFlow API

## Incoming Payloads Disabled:
1. Payload reached Runner
2. Post Payload on AlertFlow API and mark it as pending to take by any runner

### Why should I disable it?
Maybe you don't want heavy work done by this runner. Or use it just as a proxy for your Payloads and Actions should be executed by the AlertFlow Runners.

## How does the Runner report?
The Alertflow API exposes and API endpoint which catches and job progression statuses.

## Plugins
There are plans but it's complicated
