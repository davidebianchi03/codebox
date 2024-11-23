<div align="center">
    <img src="./app/src/assets/images/logo-black.png#gh-light-mode-only" style="width: 200px">
    <img src="./app/src/assets/images/logo-white.png#gh-dark-mode-only" style="width: 200px">

  <h3>
    Remote Development Environments based on Devcontainer
  </h3>
    <img alt="Docker image size" src="https://badgen.net/docker/size/dadebia/codebox?icon=docker&label=image%20size">
    <img alt="Docker image size" src="https://badgen.net/docker/pulls/dadebia/codebox?icon=docker&label=pulls">

  <br>
  <br>

</div>

**Codebox** is a service that allows developers to create remote workspaces based on docker containers. Codebox workspaces are based on [Devcontainers specification](https://containers.dev/). 

## Quickstart

The easiest way to deploy your Codebox instance is using [docker compose](./docker-compose.yml) provided in this repository.

```yaml
version: '3'

services:
  redis:
    image: redis:7.4.1
    restart: always

  codebox:
    image: dadebia/codebox:latest
    ports:
      - ${CODEBOX_PORT:-12800}:8000
    volumes:
      - codeboxdb:/codebox/db
      - codeboxdata:/codebox/data
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - CODEBOX_USE_GRAVATAR=true
      - CODEBOX_USE_SUBDOMAINS=true
      - CODEBOX_WORKSPACE_OBJECTS_PREFIX=codebox
    restart: always

volumes:
  codeboxdb:
  codeboxdata:
```

## Connect to workspace container

You can connect to a running workspace container using [**codebox-cli**](https://github.com/davidebianchi03/codebox-cli). Use official [**VS Code Extension**](https://github.com/davidebianchi03/codebox-cli) to open any workspace in VS Code with a single click.

## Disclaimer

This project is built upon devcontainers but is not affiliated with or endorsed by Microsoft or the devcontainers development team.
