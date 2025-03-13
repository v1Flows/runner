# AlertFlow Runner
> This runner is the execution part of the AlertFlow platform. Please see the [AlertFlow](https://github.com/v1Flows/AlertFlow) repo for detailed informations

## Table of Contents

- [Features](#features)
- [Configuration](#configuration)
- [Plugins](#plugins)
- [Modes](#modes)
- [Self Hosting](#self-hosting)
- [Contributing](#contributing)
- [License](#license)

## Features
- **Modes**:
- **Plugins**:
- **Docker/Kubernetes**:

## Configuration
```yaml
---

log_level: info
mode: master

alertflow:
  url: https://alertflow.org
  runner_id: null
  api_key: null

plugins:
  - name: alertmanager
    repository: https://github.com/AlertFlow/rp-alertmanager
    version: v1.0.2

alert_endpoints:
  port: 8081
```

## Plugins


## Modes

### Master
All components are enabled. The runner will receive payloads, process them and scan for pending jobs.

### Worker
The Worker mode will disable the payload receiver component. The runner will only act as an Job executor.

### Listener
The runner will only act as a payload receiver. There will be no components enable to scan or execute any jobs.

## Self Hosting

## Contributing

We welcome contributions to this project! To contribute, follow these steps:

1. Fork the repository.
2. Create a new branch:
    ```sh
    git checkout -b feature/your-feature-name
    ```
3. Make your changes and commit them:
    ```sh
    git commit -m "Add your commit message"
    ```
4. Push to the branch:
    ```sh
    git push origin feature/your-feature-name
    ```
5. Open a pull request on GitHub.

## License
This project is licensed under the GNU AFFERO GENERAL PUBLIC LICENSE Version 3. See the [LICENSE](https://github.com/AlertFlow/alertflow/blob/main/LICENSE) file for details.
