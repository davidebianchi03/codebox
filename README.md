<div align="center">
  <img src="./docs/codebox-logo.png" width="200">

  <h3>
    Remote Development Environments
  </h3>
    <img alt="Docker image size" src="https://badgen.net/docker/size/dadebia/codebox?icon=docker&label=image%20size">
    <img alt="Docker image size" src="https://badgen.net/docker/pulls/dadebia/codebox?icon=docker&label=pulls">

  <img src="./docs/preview.png" style="max-width: 550px">

  <br>
</div>

> [!WARNING]  
> This software is still an alpha version, it has many bugs.

**Codebox** is a service that allows developers to create remote workspaces. The structure of Codebox workspaces can be defined using standard spefications such as docker-compose, devcontainer, etc...



## Quickstart

The easiest way to deploy your Codebox instance is using [docker compose](./docker-compose.yml) provided in this repository.

```yaml
version: '3'

services:
  redis:
    image: redis:7.4.1
    restart: always

  db:
    image: mysql:8.0.41
    environment:
      MYSQL_ROOT_PASSWORD: ${CODEBOX_DB_ROOT_PASSWORD:-password}
      MYSQL_DATABASE: ${CODEBOX_DB_NAME:-codebox}
      MYSQL_USER: ${CODEBOX_DB_USER:-codebox}
      MYSQL_PASSWORD: ${CODEBOX_DB_PASSWORD:-password}
    healthcheck:
      test: ["CMD", "mysqladmin" ,"ping", "-h", "localhost"]
      timeout: 20s
      retries: 10
    volumes:
      - codeboxdb:/var/lib/mysql
    restart: always

  codebox:
    image: dadebia/codebox:${CODEBOX_VERSION:-latest}
    depends_on:
      db:
        condition: service_healthy
    ports:
      - ${CODEBOX_PORT:-12800}:8000
    volumes:
      - codeboxdata:/codebox/data
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - CODEBOX_EXTERNAL_URL=${CODEBOX_EXTERNAL_URL:-https://codebox.example.com}
      - CODEBOX_WILDCARD_DOMAIN=${CODEBOX_WILDCARD_DOMAIN:-codebox.example.com}
      - CODEBOX_USE_GRAVATAR=true
      - CODEBOX_USE_SUBDOMAINS=true
      - CODEBOX_DB_NAME=${CODEBOX_DB_NAME:-codebox}
      - CODEBOX_DB_USER=${CODEBOX_DB_USER:-codebox}
      - CODEBOX_DB_PASSWORD=${CODEBOX_DB_PASSWORD:-password}
    restart: always

volumes:
  codeboxdb:
  codeboxdata:
```

## Runners
Codebox cannot run workspaces by itself, you need to connect runners.

Codebox runners are services connected to the Codebox instance, they create and manage your workspaces. Each workspace type has its own runner.

Only a runner for docker-based workspaces is currently available, this runner can create workspaces using based on`docker compose` or `devcontainer`. The source code and a guide for this runner are available [here](https://gitlab.com/codebox4073715/codebox-docker-runner).

## Connect to workspace container

You can connect to a running workspace container using [**codebox-cli**](https://gitlab.com/codebox4073715/codebox-cli). Use official [**VS Code Extension**](https://gitlab.com/codebox4073715/codebox-vscode-extension) to open any workspace in VS Code with a single click.
