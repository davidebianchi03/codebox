# Server

This document describes the procedure for deploying a self-hosted instance of Codebox.

## Deploy the Server

The recommended deployment method is to use the Docker stack defined in the `docker-compose.yml` file provided in the Codebox repository.

```{warning}
Docker and Docker Compose are required. Ensure they are installed and available in your system PATH before continuing.
```

To deploy Codebox, first download the `docker-compose.yml` file:

```bash
wget [https://gitlab.com/api/v4/projects/68940432/repository/files/docker-compose.yml/raw?ref=master](https://gitlab.com/api/v4/projects/68940432/repository/files/docker-compose.yml/raw?ref=master) -O docker-compose.yml
```

Alternatively, use `curl`:

```bash
curl --output docker-compose.yml "[https://gitlab.com/api/v4/projects/68940432/repository/files/docker-compose.yml/raw?ref=master](https://gitlab.com/api/v4/projects/68940432/repository/files/docker-compose.yml/raw?ref=master)"
```

### Required Environment Variables

The Docker stack requires the following environment variables in order to start correctly:

* `CODEBOX_EXTERNAL_URL`: The externally accessible URL of the Codebox instance (for example, `https://codebox.my-domain.com`).
* `CODEBOX_WILDCARD_DOMAIN`: Codebox supports exposing HTTP-based services through dynamically generated subdomains. This value defines the wildcard domain used for this purpose. A DNS record such as `*.codebox.my-domain.com` must be configured and point to the Codebox host.

### Optional: Email Sender Configuration

Configuring an email sender service is strongly recommended. The email sender must be an SMTP-compatible server. Certain features (such as user sign-up and account approval workflows) are disabled if email delivery is not configured.

The following environment variables are used to configure the SMTP service:

* `CODEBOX_EMAIL_SMTP_HOST`: Hostname of the SMTP server.
* `CODEBOX_EMAIL_SMTP_PORT`: Port on which the SMTP server is listening.
* `CODEBOX_EMAIL_SMTP_USER`: Username used to authenticate with the SMTP server.
* `CODEBOX_EMAIL_SMTP_PASSWORD`: Password used to authenticate with the SMTP server.

A complete list of supported configuration parameters is available in the [configuration](../advanced/server-configuration) documentation.

### Starting the Stack

Once the environment variables have been defined, start the Docker stack:

```bash
docker compose up
```

```{note}
Codebox can also be deployed using Portainer by creating a new stack and pasting the contents of the `docker-compose.yml` file.
```

### Reverse Proxy Configuration

If the Codebox instance is exposed through an HTTP reverse proxy (for example, Nginx or Apache), WebSocket support **must** be enabled to ensure proper operation of interactive features.

Refer to the official documentation of your proxy server for configuration details:

* [Enable WebSockets with Nginx](https://nginx.org/en/docs/http/websocket.html)
* [Enable WebSockets with Apache](https://httpd.apache.org/docs/current/mod/mod_proxy_wstunnel.html)
