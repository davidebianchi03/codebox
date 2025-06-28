# Quickstart
This guide contains the details on the steps to install a self-hosted instance of Codebox.

## Installation
The recommended installation procedure involves using the Docker stack defined in the `docker-compose.yml` file found in the Codebox repository.

```{warning}
   Docker and Docker Compose are required, please install them before proceeding
```

To install Codebox, the first step is to download the docker-compose.yml file using the following command:

```bash
wget https://gitlab.com/api/v4/projects/68940432/repository/files/docker-compose.yml/raw?ref=master -O docker-compose.yml
```
or if you prefer to use curl:
```bash
curl --output docker-compose.yml "https://gitlab.com/api/v4/projects/68940432/repository/files/docker-compose.yml/raw?ref=master"
```

The stack requires two environment variables to start:

- `CODEBOX_EXTERNAL_URL`: the url of the codebox instance
- `CODEBOX_WILDCARD_DOMAIN`: codebox allows you to expose ports running HTTP-based services either public or with password authentication. The ports will be exposed through subdomains of this domain. You will need to define a DNS record with a name such as `*.codebox.my-domain.com`.

You can see the list of all available parameters [here](./configuration.md).

Now you can start your docker stack:
```
docker compose up
```

```{note}
You can also launch your Codebox instance using Portainer by creating a new stack and copying the contents of the `docker-compose.yml` file.
```

## Proxy
If you are exposing your Codebox instance through an HTTP proxy, such as Nginx or Apache, you will need to enable WebSocket connections.
