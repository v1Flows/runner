# Runner
> This runner is the execution engine of the AlertFlow and exFlow platform. Please see the [AlertFlow](https://github.com/v1Flows/AlertFlow) or [exFlow](https://github.com/v1Flows/exFlow) repo for detailed informations

## Table of Contents

- [Features](#features)
- [Configuration](#configuration)
- [Plugins](#plugins)
- [Modes](#modes)
- [Self Hosting](#self-hosting)
- [Contributing](#contributing)
- [License](#license)

## Features
- **Modes**: The runner can be started in different modes which either offer full functionality or just be a standby listener for incoming alerts
- **Plugins**: Develop your own plugins or use our existing ones to extend the functionality of this runner and alertflow / exflow to your needs

## Configuration

To conntect an runner to exFlow or AlertFlow you first have to set them up and copy the runner_id and or the api key from the created project. As an Admin you can copy the Global Share Runner token from the admin view.

```yaml
---

log_level: info
mode: master

alertflow:
  enabled: true
  url: https://alertflow.org
  runner_id: null
  api_key: null

exflow:
  enabled: true
  url: https://exflow.org
  runner_id: null
  api_key: null

plugins:
  - name: alertmanager
    version: v1.2.4
  - name: git
    version: v1.2.0
  - name: ansible
    version: v1.3.2
  - name: ssh
    version: v1.4.0

api_endpoint:
  port: 8081
```

## Plugins
The runner can be extended by integrating plugins following a specific schema. A list of available plugins can be found [here](https://github.com/v1Flows/runner-plugins).

To develop your own plugin you can start right away with this [template](https://github.com/v1Flows/runner-plugins/tree/develop/template)

## Modes

### Master
All components are enabled. The runner will receive payloads, process them and scan for pending jobs.

### Worker
The Worker mode will disable the payload receiver component. The runner will only act as an Job executor.

### Listener
The runner will only act as a payload receiver. There will be no components enable to scan or execute any jobs.

## Self Hosting
To host the Runner on your own infrastructure we provide various docker images available at 
[Docker Hub](htthttps://hub.docker.com/r/justnz/runner).
- **justnz/runner:latest** - Latest Version
- **justnz/runner:vx.x.x** - Versioned release

```sh
docker run -p 8081:8081 -v /your/config/path/config.yaml:/app/config/config.yaml justnz/runner:latest
```

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
This project is licensed under the GNU AFFERO GENERAL PUBLIC LICENSE Version 3. See the [LICENSE](https://github.com/v1Flows/runner/blob/main/LICENSE) file for details.
